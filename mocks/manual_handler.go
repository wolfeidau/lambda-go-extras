// Package mocks contains manual mock implementations.
package mocks

import (
	"context"
)

// Handler is a manual mock implementation of lambda.Handler.
type Handler struct {
	InvokeFunc func(ctx context.Context, payload []byte) ([]byte, error)
}

// Invoke calls the injected InvokeFunc if set.
func (m *Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	if m.InvokeFunc != nil {
		return m.InvokeFunc(ctx, payload)
	}
	return nil, nil
}
