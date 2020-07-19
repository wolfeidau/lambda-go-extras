package raw

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
)

// Option assign settings to the zerolog middlware
type Option func(opts *rawOptions)

// settings for the zerolog middlware
type rawOptions struct {
	fields map[string]interface{}
	output io.Writer
}

// Fields pass a map of attributes which are appended to all log messages
// emitted by this logger.
func Fields(fields map[string]interface{}) Option {
	return func(opts *rawOptions) {
		for k, v := range fields {
			opts.fields[k] = v
		}
	}
}

// Output is a writer where logs in JSON format are written.
// Defaults to os.Stderr.
func Output(output io.Writer) Option {
	return func(opts *rawOptions) {
		opts.output = output
	}
}

// New build a new raw event logging middleware, this uses zerolog to emit
// a log message for the input and output events
func New(options ...Option) func(next lambda.Handler) lambda.Handler {

	opts := &rawOptions{
		output: os.Stderr,
		fields: make(map[string]interface{}),
	}

	for _, opt := range options {
		opt(opts)
	}

	return func(next lambda.Handler) lambda.Handler {
		return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			lc, _ := lambdacontext.FromContext(ctx)

			zlog := zerolog.New(opts.output).With().
				Stack().
				Fields(opts.fields).
				Str("aws_request_id", lc.AwsRequestID).
				Str("amzn_trace_id", os.Getenv("_X_AMZN_TRACE_ID")).
				Logger()

			var v interface{}

			if ok := unmarshal(payload, &v); ok {
				zlog.Info().Fields(map[string]interface{}{
					"event": v,
				}).Msg("incoming event")
			}

			result, err := next.Invoke(ctx, payload)
			if err != nil {
				return nil, err
			}

			if ok := unmarshal(result, &v); ok {
				zlog.Info().Err(err).Fields(map[string]interface{}{
					"event": v,
				}).Msg("outgoing event")
			}

			return result, err
		})
	}
}

func unmarshal(payload []byte, v interface{}) bool {
	err := json.Unmarshal(payload, &v)
	return err == nil
}
