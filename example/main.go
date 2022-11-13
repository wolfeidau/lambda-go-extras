package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

var (
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	flds := lmw.FieldMap{"commit": commit, "buildDate": buildDate}

	ch := lmw.New(
		raw.New(raw.Fields(flds)),   // raw event logger which prints input and output of handler
		zlog.New(zlog.Fields(flds)), // inject zerolog into the context
	).Then(lambdaextras.GenericHandler(processSQSEvent))

	// use StartWithOptions as StartHandler is deprecated
	lambda.StartWithOptions(ch)
}

func processSQSEvent(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {
	zerolog.Ctx(ctx).Info().Msg("sqsEvent")
	return "ok", nil
}
