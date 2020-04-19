# lambda-go-extras

This module provides a middleware layer for [github.com/aws/aws-lambda-go](https://github.com/aws/aws-lambda-go). This project is heavily based on https://github.com/justinas/alice.

[![GitHub Actions status](https://github.com/wolfeidau/lambda-go-extras/workflows/Go/badge.svg?branch=master)](https://github.com/wolfeidau/lambda-go-extras/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/wolfeidau/lambda-go-extras)](https://goreportcard.com/report/github.com/wolfeidau/lambda-go-extras)
[![Documentation](https://godoc.org/github.com/wolfeidau/lambda-go-extras?status.svg)](https://godoc.org/github.com/wolfeidau/lambda-go-extras) [![Coverage Status](https://coveralls.io/repos/github/wolfeidau/lambda-go-extras/badge.svg?branch=master)](https://coveralls.io/github/wolfeidau/lambda-go-extras?branch=master)

# Usage

This illustrates how easy it is to dump the input and output payloads using a simple middleware chain.

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
)

func main() {

	ch := lmw.New(logEvent).ThenFunc(processEvent)

	lambda.StartHandler(ch)
}

func logEvent(next lambda.Handler) lambda.Handler {
	return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {

		fmt.Println(string(payload))

		result, err := next.Invoke(ctx, payload)

		fmt.Println(string(result))

		return result, err
	})
}

func processEvent(ctx context.Context, payload []byte) ([]byte, error) {
	return []byte("ok"), nil
}
```

# License

This code was authored by [Mark Wolfe](https://www.wolfe.id.au) and is licensed under the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0).