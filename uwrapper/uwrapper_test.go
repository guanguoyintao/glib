package uwrapper

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGoExecHandlerWithTimeout(t *testing.T) {
	type args struct {
		ctx     context.Context
		f       wrapperFunc
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "timeout",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context) error {
					fmt.Println("-------timeout-----start------------")
					time.Sleep(10 * time.Minute)
					fmt.Println("-------timeout-----end------------")

					return nil
				},
				timeout: 2 * time.Second,
			},
		},
		{
			name: "normal return",
			args: args{
				ctx: context.Background(),
				f: func(ctx context.Context) error {
					fmt.Println("-------normal return-----start------------")
					time.Sleep(1 * time.Second)
					fmt.Println("-------normal return-----end------------")

					return nil
				},
				timeout: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GoExecHandlerWithTimeout(tt.args.ctx, tt.args.f, tt.args.timeout)
		})
	}
	select {}
}
