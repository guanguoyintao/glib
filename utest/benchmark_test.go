package utest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	type args struct {
		ctx            context.Context
		concurrentList []uint32
	}
	tests := []struct {
		name string
		f    TestFunc
		args args
	}{
		{
			name: "1",
			f: func(ctx context.Context) error {
				i := rand.Intn(10)
				time.Sleep(time.Duration(i) * time.Second)
				if i == 5 || i == 0 {
					return fmt.Errorf("err")
				}
				return nil
			},
			args: args{
				ctx:            context.Background(),
				concurrentList: []uint32{100, 1000, 10000},
			},
		},
		{
			name: "2",
			f: func(ctx context.Context) error {
				time.Sleep(2 * time.Second)
				return nil
			},
			args: args{
				ctx:            context.Background(),
				concurrentList: []uint32{100, 1000, 10000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.TestClient(tt.args.ctx, tt.args.concurrentList)
		})
	}
}
