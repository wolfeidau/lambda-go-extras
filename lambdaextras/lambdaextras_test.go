package lambdaextras

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerFunc(t *testing.T) {
	assert := require.New(t)

	h := HandlerFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		return bytes.Replace(payload, []byte("hello"), []byte("world"), 1), nil
	})

	data, err := h.Invoke(context.TODO(), []byte("hello"))
	assert.NoError(err)
	assert.Equal([]byte("world"), data)
}
