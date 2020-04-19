package lambdaextras

import "context"

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as Lambda Handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ctx context.Context, payload []byte) ([]byte, error)

// Invoke calls f(ctx, payload).
func (f HandlerFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return f(ctx, payload)
}
