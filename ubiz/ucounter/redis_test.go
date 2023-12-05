package ucounter

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRedisCounter(t *testing.T) {
	cronJobInterval = 20 * time.Second
	type counterConfig struct {
		messageID uint64
		increment uint32
		command   commandType
		want      int32
		wantErr   bool
		wait      bool
		timeSleep time.Duration
	}
	type args struct {
		ctx     context.Context
		key     CounterNameSpaceType
		configs []counterConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				ctx: context.Background(),
				key: "ucounter-test",
				configs: []counterConfig{
					{
						messageID: 0,
						increment: 0,
						command:   getCommand,
						want:      0,
						wait:      true,
					},
					{
						messageID: 0,
						increment: 2,
						command:   incrCommand,
						want:      -1,
						wantErr:   false,
						wait:      false,
					},
					{
						messageID: 0,
						increment: 2,
						command:   incrCommand,
						want:      -1,
						wantErr:   false,
						wait:      false,
					},
					{
						messageID: 0,
						increment: 0,
						command:   getCommand,
						want:      4,
						wait:      true,
						timeSleep: cronJobInterval,
					},
					{
						messageID: 0,
						increment: 2,
						command:   incrCommand,
						want:      -1,
						wantErr:   false,
						wait:      false,
					},
					{
						messageID: 0,
						increment: 2,
						command:   incrCommand,
						want:      -1,
						wantErr:   false,
						wait:      false,
					},
					{
						messageID: 0,
						increment: 0,
						command:   getCommand,
						want:      8,
						wantErr:   false,
						wait:      true,
					},
				},
			},
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, conf := range tt.args.configs {
				counter, err := NewCounter(tt.args.ctx, tt.args.key)
				if (err != nil) != conf.wantErr {
					t.Error(err)
				}
				fn := func() {
					defer wg.Done()
					switch conf.command {
					case incrCommand:
						res, err := counter.Incr(tt.args.ctx, conf.messageID, conf.increment)
						if err != nil {
							t.Error(err)
						}
						fmt.Printf("msg:%d,incr:%d,incr res:%d", conf.messageID, conf.increment, res)
					case decrCommand:
						res, err := counter.Decr(tt.args.ctx, conf.messageID, conf.increment)
						if err != nil {
							t.Error(err)
						}
						fmt.Printf("msg:%d,incr:%d,decr res:%d", conf.messageID, conf.increment, res)
					case getCommand:
						counterNum, err := counter.Get(tt.args.ctx, conf.messageID)
						if err != nil {
							t.Error(err)
						}
						fmt.Printf("counter num is %v\n", counterNum)
						assert.Equal(t, counterNum, uint32(conf.want))
					}
				}
				wg.Wait()
				wg.Add(1)
				if conf.wait {
					fn()
				} else {
					go fn()
				}
				time.Sleep(conf.timeSleep)
			}
		})
	}
}

func TestRedisCounter_cronJob(t *testing.T) {
	type args struct {
		ctx context.Context
		key CounterNameSpaceType
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				ctx: context.Background(),
				key: "ucounter-test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCounter(tt.args.ctx, tt.args.key)
			if err != nil {
				t.Error(err)
			}
			select {}
		})
	}
}

func TestRedisCounter_Incr(t *testing.T) {
	type args struct {
		ctx    context.Context
		key    CounterNameSpaceType
		megID  uint64
		number uint32
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				ctx:    context.Background(),
				key:    "ucounter-test",
				megID:  200,
				number: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, err := NewCounter(tt.args.ctx, tt.args.key)
			if err != nil {
				t.Error(err)
			}
			res, err := counter.Incr(tt.args.ctx, tt.args.megID, tt.args.number)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(res)
		})
	}
}

func TestRedisCounter_Get(t *testing.T) {
	type args struct {
		ctx   context.Context
		key   CounterNameSpaceType
		megID uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				ctx:   context.Background(),
				key:   "ucounter-test",
				megID: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, err := NewCounter(tt.args.ctx, tt.args.key)
			if err != nil {
				t.Error(err)
			}
			number, err := counter.Get(tt.args.ctx, tt.args.megID)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(number)
		})
		select {}
	}
}
