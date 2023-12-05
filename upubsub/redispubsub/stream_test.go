package redispubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/upubsub"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/config"
	"github.com/go-playground/assert/v2"
	"os"
	"path"
	"sync"
	"testing"
)

var (
	testStreamMQ upubsub.MQ
)

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(currentPath)
	config.Init(path.Join(currentPath, "conf"))
	fmt.Printf("config init %+v", config.GetConfig())
	cache.Init(config.GetConfig())
	testStreamMQ, err = NewRedisStreamMQ(ctx)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	os.Exit(code)
}

func Test_redisStreamConsumer_Consume(t *testing.T) {
	type fields struct {
		topics string
		group  string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "group 1 topic 1 consumer 1",
			fields: fields{
				topics: "test:stream",
				group:  "1",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "group 1 topic 1 consumer 2",
			fields: fields{
				topics: "test:stream",
				group:  "1",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	mq := testStreamMQ
	consumers := make([]upubsub.Consumer, 0, len(tests))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := mq.NewConsumer(tt.args.ctx, tt.fields.group, tt.fields.topics)
			if err != nil {
				t.Error(err)
			}
			consumers = append(consumers, c)
		})
	}
	wg := sync.WaitGroup{}
	for _, consumer := range consumers {
		wg.Add(1)
		c := consumer
		go func() {
			defer wg.Done()
			msg, err := c.Consume(context.Background())
			if err != nil {
				t.Error(err)
			}
			fmt.Println("===============start===============")
			fmt.Printf("msg:%+v\n", string(msg.Message(context.Background())))
			fmt.Println("==============end================")
		}()
	}
	wg.Wait()
}

func Test_redisRedisProducer_Publish(t *testing.T) {
	type fields struct {
		topics []string
	}
	type args struct {
		ctx context.Context
		msg upubsub.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "topic 1 producer 1",
			fields: fields{
				topics: []string{"1"},
			},
			args: args{
				ctx: context.Background(),
				msg: NewRedisStreamEvent("", []byte("t1m1")),
			},
			wantErr: false,
		},
		{
			name: "topic 1 producer 2",
			fields: fields{
				topics: []string{"1"},
			},
			args: args{
				ctx: context.Background(),
				msg: NewRedisStreamEvent("", []byte("t1m2")),
			},
			wantErr: false,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				wg.Add(1)
				defer wg.Done()
				p, err := testStreamMQ.NewProducer(tt.args.ctx, tt.fields.topics...)
				if err != nil {
					t.Error(err)
				}
				err = p.Publish(tt.args.ctx, tt.args.msg)
				if (err != nil) != tt.wantErr {
					t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
	wg.Wait()
}

func Test_redisRedisProducer_Publish_And_Consume(t *testing.T) {
	type fields struct {
		group  string
		topics []string
	}
	type args struct {
		ctx context.Context
		msg upubsub.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "topic 1 producer 1",
			fields: fields{
				group:  "",
				topics: []string{"test:stream"},
			},
			args: args{
				ctx: context.Background(),
				msg: NewRedisStreamEvent("", []byte("t1m3")),
			},
			wantErr: false,
		},
		{
			name: "topic 1 producer 2",
			fields: fields{
				group:  "",
				topics: []string{"test:stream"},
			},
			args: args{
				ctx: context.Background(),
				msg: NewRedisStreamEvent("", []byte("t1m4")),
			},
			wantErr: false,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				wg.Add(1)
				defer wg.Done()
				p, err := testStreamMQ.NewProducer(tt.args.ctx, tt.fields.topics...)
				if err != nil {
					t.Error(err)
				}
				err = p.Publish(tt.args.ctx, tt.args.msg)
				if (err != nil) != tt.wantErr {
					t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				}
				consumer, err := testStreamMQ.NewConsumer(tt.args.ctx, tt.fields.group, tt.fields.topics[0])
				if err != nil {
					t.Error(err)
				}
				message, err := consumer.Consume(tt.args.ctx)
				if err != nil {
					t.Error(err)
				}
				fmt.Println(string(message.Message(tt.args.ctx)))
				assert.IsEqual(message.Message(tt.args.ctx), tt.args.msg)
				err = message.Ack(tt.args.ctx)
				if err != nil {
					t.Error(err)
				}
			}()
		})
	}
	wg.Wait()
}

func Test_Consume(t *testing.T) {
	type fields struct {
		group  string
		topics []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ushow:result_queue:stream",
			fields: fields{
				group:  "test",
				topics: []string{"ushow:result_queue:stream"},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	wg := sync.WaitGroup{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				wg.Add(1)
				defer wg.Done()
				consumer, err := testStreamMQ.NewConsumer(tt.args.ctx, tt.fields.group, tt.fields.topics[0])
				if err != nil {
					t.Error(err)
				}
				message, err := consumer.Consume(tt.args.ctx)
				if err != nil {
					t.Error(err)
				}
				fmt.Println(string(message.Message(tt.args.ctx)))
			}()
		})
	}
	wg.Wait()
}
