package middleware

import (
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	"github.com/wolfeidau/lambda-go-extras/mocks"
)

func TestNew(t *testing.T) {

	assert := require.New(t)

	c1 := func(h lambda.Handler) lambda.Handler {
		return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			return h.Invoke(ctx, payload)
		})
	}

	middlewares := []Middleware{c1}

	chain := New(middlewares...)

	assert.Len(chain.middlewares, 1)
}

func TestThen(t *testing.T) {

	assert := require.New(t)

	c1 := func(h lambda.Handler) lambda.Handler {
		return lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			return h.Invoke(ctx, payload)
		})
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := mocks.NewMockHandler(ctrl)

	handler.EXPECT().Invoke(gomock.Any(), []byte("hello")).Return([]byte("world"), nil)

	ch := New(c1).Then(handler)

	data, err := ch.Invoke(context.TODO(), []byte("hello"))
	assert.NoError(err)
	assert.Equal([]byte("world"), data)
}

func TestThenFunc(t *testing.T) {
	assert := require.New(t)

	ch := New().ThenFunc(lambdaextras.HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		return bytes.Replace(payload, []byte("hello"), []byte("world"), 1), nil
	}))

	data, err := ch.Invoke(context.TODO(), []byte("hello"))
	assert.NoError(err)
	assert.Equal([]byte("world"), data)

}
