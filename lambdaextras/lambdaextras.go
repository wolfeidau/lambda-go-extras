// Package lambdaextras contains extras for building Go based lambdas.
package lambdaextras

import (
	"context"
	"encoding/json"
)

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as Lambda Handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ctx context.Context, payload []byte) ([]byte, error)

// Invoke calls f(ctx, payload).
func (f HandlerFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return f(ctx, payload)
}

// GenericHandler converts a typed lambda handler into a Handler func.
func GenericHandler[I any, O any](handler func(context.Context, I) (O, error)) HandlerFunc {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var input I

		err := json.Unmarshal(payload, &input)
		if err != nil {
			return nil, err
		}

		output, err := handler(ctx, input)
		if err != nil {
			return nil, err
		}

		return json.Marshal(&output)
	}
}
