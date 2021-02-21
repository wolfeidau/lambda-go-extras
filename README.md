# lambda-go-extras

This module provides a middleware layer for [github.com/aws/aws-lambda-go](https://github.com/aws/aws-lambda-go). This project is heavily based on https://github.com/justinas/alice.

[![GitHub Actions status](https://github.com/wolfeidau/lambda-go-extras/workflows/Go/badge.svg?branch=master)](https://github.com/wolfeidau/lambda-go-extras/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/lambda-go-extras)](https://goreportcard.com/report/github.com/wolfeidau/lambda-go-extras)
[![Documentation](https://godoc.org/github.com/wolfeidau/lambda-go-extras?status.svg)](https://godoc.org/github.com/wolfeidau/lambda-go-extras)

# Why?

Having used the `github.com/aws/aws-lambda-go` package for a while now I have found it annoying switching between the Go standard libraries `http` package and this library. After some review I think the thing I miss the most is the ability to chain a list of handlers, with each link responsible for a part of the puzzle.

Being able to compose these chains offers a lot of flexibility and reuse across projects.

Given the default way of using the `github.com/aws/aws-lambda-go` is via the [Start(handler interface{})](https://godoc.org/github.com/aws/aws-lambda-go/lambda#Start) function, which is very flexible, but not easily extended due it's dynamic nature. So I have moved to using the less used [func StartHandler(handler Handler)](https://godoc.org/github.com/aws/aws-lambda-go/lambda#StartHandler) which is more idiomatic with its `Handler` being just an interface with a single `Invoke` method. This module has extended this with a simple middleware chain inspired by alice as mentioned above.

# Usage

This illustrates how easy it is to dump the input and output payloads using a simple middleware chain.

```go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

// assigned during build time with -ldflags
var (
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {

	flds := lmw.FieldMap{"commit": commit, "buildDate": buildDate}

	ch := lmw.New(
		raw.New(raw.Fields(flds)),   // raw event logger primarily used during development
		zlog.New(zlog.Fields(flds)), // inject zerolog into the context
	).ThenFunc(processEvent)

	lambda.StartHandler(ch)
}

func processEvent(ctx context.Context, payload []byte) ([]byte, error) {
	zerolog.Ctx(ctx).Info().Msg("processEvent")
	return []byte("ok"), nil
}
```

There are two bundled middleware, these are 

* `github.com/wolfeidau/lambda-go-extras/middleware/raw` - raw event logger which outputs input and output events which also includes `aws_request_id` and the fields you provide at initialisation.
* `github.com/wolfeidau/lambda-go-extras/middleware/zerolog` - configures the context with a [zerolog](https://github.com/rs/zerolog) logger which includes `aws_request_id` and the fields you provide at initialisation.

# License

This code was authored by [Mark Wolfe](https://www.wolfe.id.au) and is licensed under the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0).
