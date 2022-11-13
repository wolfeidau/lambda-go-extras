package standard

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	"github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

type defaultConfig struct {
	rawEnabled bool
	fields     middleware.FieldMap
	output     io.Writer
	level      zerolog.Level
}

type Option func(config *defaultConfig)

// RawEnabled is a flag to toggle this middleware on or off.
// Defaults to true.
func RawEnabled(flag bool) Option {
	return func(config *defaultConfig) {
		config.rawEnabled = flag
	}
}

// Output is a writer where logs in JSON format are written.
// Defaults to os.Stderr.
func Output(output io.Writer) Option {
	return func(config *defaultConfig) {
		config.output = output
	}
}

// LogLevel minimum accepted level for logging.
// Defaults to zerolog.InfoLevel.
func LogLevel(level zerolog.Level) Option {
	return func(opts *defaultConfig) {
		opts.level = level
	}
}

// Default configures a standard lambda handler with default middleware which includes the zerolog and raw logging.
func Default(h lambda.Handler, opts ...Option) {

	config := &defaultConfig{
		rawEnabled: true,
		fields:     make(middleware.FieldMap),
		output:     os.Stderr,
		level:      zerolog.InfoLevel,
	}

	for _, opt := range opts {
		opt(config)
	}

	ch := middleware.New(
		raw.New(raw.Fields(config.fields), raw.Enabled(config.rawEnabled), raw.Output(config.output)),
		zlog.New(zlog.Fields(config.fields), zlog.Output(config.output)),
	)

	// use StartWithOptions as StartHandler is deprecated
	lambda.StartWithOptions(ch.Then(h))
}

// GenericDefault configures a generic standard lambda handler with default middleware which includes the zerolog and raw logging.
func GenericDefault[I any, O any](handler func(context.Context, I) (O, error), opts ...Option) {
	Default(lambdaextras.GenericHandler(handler), opts...)
}
