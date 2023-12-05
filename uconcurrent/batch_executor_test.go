package uconcurrent

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatchExecutor_Run(t *testing.T) {
	contextKey := "k"
	contextValue := "v"

	type fields struct {
		HasTimeout bool
		runnables  []Runnable
		results    []interface{}
	}
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "all success",
			fields: fields{
				HasTimeout: false,
				runnables: []Runnable{
					func(ctx context.Context) (interface{}, error) {
						// 确保context正确传递
						//assert.Equal(t, contextValue, ctx.Value(contextKey))
						time.Sleep(1 * time.Second)
						return 1, nil
					},
					func(ctx context.Context) (interface{}, error) {
						// 确保context正确传递
						//assert.Equal(t, contextValue, ctx.Value(contextKey))
						time.Sleep(1 * time.Second)
						return 2, nil
					},
					func(ctx context.Context) (interface{}, error) {
						// 确保context正确传递
						//assert.Equal(t, contextValue, ctx.Value(contextKey))
						time.Sleep(1 * time.Second)
						return 3, nil
					},
				},
				results: []interface{}{1, 2, 3},
			},
			args: args{timeout: 2 * time.Second},
		},
		{
			name: "one timeout",
			fields: fields{
				HasTimeout: true,
				runnables: []Runnable{
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(900 * time.Millisecond)
						return 1, nil
					},
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(2 * time.Second)
						return 2, nil
					},
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(900 * time.Millisecond)
						return 3, nil
					},
				},
				results: []interface{}{1, nil, 3},
			},
			args: args{timeout: 1 * time.Second},
		},
		{
			name: "all timeout",
			fields: fields{
				HasTimeout: true,
				runnables: []Runnable{
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(900 * time.Millisecond)
						return 1, nil
					},
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(2 * time.Second)
						return 2, nil
					},
					func(ctx context.Context) (interface{}, error) {
						time.Sleep(900 * time.Millisecond)
						return 3, nil
					},
				},
				results: []interface{}{nil, nil, nil},
			},
			args: args{timeout: 100 * time.Millisecond},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BatchExecutor{}
			for _, runnable := range tt.fields.runnables {
				b.AddRunnable(runnable)
			}
			ctx := context.WithValue(context.Background(), contextKey, contextValue)
			b.Run(ctx, tt.args.timeout)
			for i, result := range b.Results {
				assert.Equal(t, tt.fields.results[i], result)
			}
			fmt.Println(b.Results, b.SuccessCount, b.ErrorCount)
			assert.Equal(t, tt.fields.HasTimeout, b.HasTimeout)
		})
	}
}
