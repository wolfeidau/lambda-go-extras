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
	rawEnabled  bool
	fields      middleware.FieldMap
	output      io.Writer
	level       zerolog.Level
	middlewares []middleware.Middleware
}

type Option func(config *defaultConfig)

// Fields pass a map of attributes which are appended to all log messages and raw events.
func Fields(fields map[string]interface{}) Option {
	return func(opts *defaultConfig) {
		for k, v := range fields {
			opts.fields[k] = v
		}
	}
}

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

// Append middleware to the chain.
// Defaults to raw and zerolog middleware.
func Append(mw middleware.Middleware) Option {
	return func(opts *defaultConfig) {
		opts.middlewares = append(opts.middlewares, mw)
	}
}

// Default configures a standard lambda handler with default middleware which includes the zerolog and raw logging.
func Default(h lambda.Handler, opts ...Option) {

	config := &defaultConfig{
		rawEnabled:  true,
		fields:      make(middleware.FieldMap),
		output:      os.Stderr,
		level:       zerolog.InfoLevel,
		middlewares: make([]middleware.Middleware, 0),
	}

	for _, opt := range opts {
		opt(config)
	}

	// prepend the two default middleware to the list, then add the ones from the configuration while preserving order
	//
	// this is done to ensure the logging middleware is at the top of stack to enable debugging of other middleware lower down the chain
	// as the zerolog logger will be in the context.
	config.middlewares = append(
		[]middleware.Middleware{
			raw.New(raw.Fields(config.fields), raw.Enabled(config.rawEnabled), raw.Output(config.output)),
			zlog.New(zlog.Fields(config.fields), zlog.Output(config.output)),
		},
		config.middlewares...,
	)

	ch := middleware.New(
		config.middlewares...,
	)

	// use StartWithOptions as StartHandler is deprecated
	lambda.StartWithOptions(ch.Then(h))
}

// GenericDefault configures a generic standard lambda handler with default middleware which includes the zerolog and raw logging.
func GenericDefault[I any, O any](handler func(context.Context, I) (O, error), opts ...Option) {
	Default(lambdaextras.GenericHandler(handler), opts...)
}
