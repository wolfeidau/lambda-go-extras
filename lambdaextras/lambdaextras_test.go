package lambdaextras

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerFunc(t *testing.T) {
	assert := require.New(t)

	h := HandlerFunc(func(_ context.Context, payload []byte) ([]byte, error) {
		return bytes.Replace(payload, []byte("hello"), []byte("world"), 1), nil
	})

	data, err := h.Invoke(context.TODO(), []byte("hello"))
	assert.NoError(err)
	assert.Equal([]byte("world"), data)
}

type In struct {
	Msg string `json:"msg,omitempty"`
}

type Out struct {
	Result string `json:"result,omitempty"`
}

func TestGenericHandler(t *testing.T) {
	assert := require.New(t)

	h := GenericHandler(func(_ context.Context, input In) (Out, error) {
		return Out{Result: input.Msg}, nil
	})

	data, err := h.Invoke(context.TODO(), []byte(`{"msg": "hello"}`))
	assert.NoError(err)
	assert.Equal([]byte(`{"result":"hello"}`), data)
}
