package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
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
	).ThenFunc(processEvent)

	lambda.StartHandler(ch)
}

func processEvent(ctx context.Context, payload []byte) ([]byte, error) {
	zerolog.Ctx(ctx).Info().Msg("processEvent")
	return []byte("ok"), nil
}
