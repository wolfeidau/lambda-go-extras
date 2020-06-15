package zlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/lambda-go-extras/middleware"
)

func testHandler(ctx context.Context, payload []byte) ([]byte, error) {
	// retrieve a logger from the context
	log := zerolog.Ctx(ctx)

	log.Info().Str("msg", "hello").Msg("testing")

	return []byte{}, nil
}

func TestNew(t *testing.T) {
	assert := require.New(t)

	_, filename, _, _ := runtime.Caller(0)
	t.Logf("Current test filename: %s", filename)

	type args struct {
		fields map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "should dump json",
			args: args{fields: FieldMap{"msg": "hello"}},
			want: map[string]string{
				"amzn_trace_id":  "",
				"aws_request_id": "test123",
				"caller":         fmt.Sprintf("%s:21", filename),
				"level":          "info",
				"message":        "testing",
				"msg":            "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := lambdacontext.NewContext(context.TODO(), &lambdacontext.LambdaContext{
				AwsRequestID: "test123",
			})

			buf := new(bytes.Buffer)
			ch := middleware.New(New(Fields(tt.args.fields), Output(buf))).ThenFunc(testHandler)

			data, err := ch.Invoke(ctx, []byte{})
			assert.NoError(err)
			assert.Equal([]byte{}, data)

			jsonWant, err := json.Marshal(&tt.want)
			assert.NoError(err)

			assert.JSONEq(string(jsonWant), buf.String())

		})
	}
}
