package zerolog

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
)

// Option assign settings to the zerolog middlware
type Option func(opts *zerlogOptions)

// settings for the zerolog middlware
type zerlogOptions struct {
	fields map[string]interface{}
	output io.Writer
}

// Fields pass a map of attributes which are appended to all log messages
// emitted by this logger.
func Fields(fields map[string]interface{}) Option {
	return func(opts *zerlogOptions) {
		for k, v := range fields {
			opts.fields[k] = v
		}
	}
}

// Output is a writer where logs in JSON format are written.
// Defaults to os.Stderr.
func Output(output io.Writer) Option {
	return func(opts *zerlogOptions) {
		opts.output = output
	}
}

// New build a new zerlog middleware with the provided configuration which has
// Stack and Caller enabled
func New(options ...Option) func(next lambda.Handler) lambda.Handler {
	opts := &zerlogOptions{
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
				Stack().Caller().
				Fields(opts.fields).
				Str("aws_request_id", lc.AwsRequestID).
				Str("amzn_trace_id", os.Getenv("_X_AMZN_TRACE_ID")).
				Logger()

			return next.Invoke(zlog.WithContext(ctx), payload)
		})
	}
}
