package raw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/lambda-go-extras/lambdaextras"
	"github.com/wolfeidau/lambda-go-extras/middleware"
)

func okHandler(_ context.Context, payload []byte) ([]byte, error) {
	return payload, nil
}

func errHandler(_ context.Context, _ []byte) ([]byte, error) {
	return nil, errors.New("boom")
}

func badHandler(_ context.Context, _ []byte) ([]byte, error) {
	return []byte("hello"), nil
}

func TestNew(t *testing.T) {
	assert := require.New(t)

	type args struct {
		fields  map[string]any
		enabled bool
	}
	tests := []struct {
		name        string
		args        args
		payload     []byte
		wantCount   int
		wantOutput  []map[string]any
		handlerFunc lambdaextras.HandlerFunc
		wantErr     bool
	}{
		{
			name:    "should dump json",
			args:    args{fields: map[string]any{"msg": "hello"}, enabled: true},
			payload: []byte(`{"source": "welcome"}`),
			wantOutput: []map[string]any{
				{
					"amzn_trace_id":  "",
					"aws_request_id": "test123",
					"event":          map[string]any{"source": "welcome"},
					"level":          "info",
					"message":        "incoming event",
					"msg":            "hello",
				}, {
					"amzn_trace_id":  "",
					"aws_request_id": "test123",
					"event":          map[string]any{"msg": "ok"},
					"level":          "info",
					"message":        "outgoing event",
					"msg":            "hello",
				},
			},
			wantCount:   3,
			handlerFunc: okHandler,
		},
		{
			name:    "should dump json with error",
			args:    args{fields: map[string]any{"msg": "hello"}, enabled: true},
			payload: []byte(`{"source": "welcome"}`),
			wantOutput: []map[string]any{
				{
					"amzn_trace_id":  "",
					"aws_request_id": "test123",
					"event":          map[string]any{"source": "welcome"},
					"level":          "info",
					"message":        "incoming event",
					"msg":            "hello",
				},
			},
			wantCount:   2,
			handlerFunc: errHandler,
			wantErr:     true,
		},
		{
			name:    "should dump json with result warning",
			args:    args{fields: map[string]any{"msg": "hello"}, enabled: true},
			payload: []byte(`{"source": "welcome"}`),
			wantOutput: []map[string]any{
				{
					"amzn_trace_id":  "",
					"aws_request_id": "test123",
					"event":          map[string]any{"source": "welcome"},
					"level":          "info",
					"message":        "incoming event",
					"msg":            "hello",
				},
			},
			wantCount:   2,
			handlerFunc: badHandler,
		},
		{
			name:    "should dump json with payload warning",
			args:    args{fields: map[string]any{"msg": "hello"}, enabled: true},
			payload: []byte(`hello`),
			wantOutput: []map[string]any{
				{
					"amzn_trace_id":  "",
					"aws_request_id": "test123",
					"event":          map[string]any{"msg": "ok"},
					"level":          "info",
					"message":        "outgoing event",
					"msg":            "hello",
				},
			},
			wantCount:   2,
			handlerFunc: okHandler,
		},
		{
			name:        "should be disabled",
			args:        args{fields: map[string]any{"msg": "hello"}, enabled: false},
			payload:     []byte(`{"source": "welcome"}`),
			wantOutput:  []map[string]any{},
			wantCount:   1,
			handlerFunc: okHandler,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			ctx := lambdacontext.NewContext(context.TODO(), &lambdacontext.LambdaContext{
				AwsRequestID: "test123",
			})

			buf := new(bytes.Buffer)
			ch := middleware.New(New(Fields(tt.args.fields), Output(buf), Enabled(tt.args.enabled))).ThenFunc(tt.handlerFunc)

			_, err := ch.Invoke(ctx, tt.payload)
			if tt.wantErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			fmt.Println(tt.name, buf.String())

			lines := strings.Split(buf.String(), "\n")
			assert.Len(lines, tt.wantCount)

			for n := range tt.wantOutput {
				jsonOutput, err := json.Marshal(&tt.wantOutput[n])
				assert.NoError(err)
				assert.JSONEq(string(jsonOutput), lines[n])
			}
		})
	}
}
