package redispubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/upubsub"
	"sync"
	"testing"
)

var (
	testPubSubMQ *redisPubSubMQ
)

//
//func TestMain(m *testing.M) {
//	var err error
//	ctx := context.Background()
//	currentPath, err := os.Getwd()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(currentPath)
//	config.Init(path.Join(currentPath, "conf"))
//	fmt.Printf("config init %+v", config.GetConfig())
//	cache.Init(config.GetConfig())
//	testPubSubMQ, err = newRedisPubSubMQ(ctx)
//	if err != nil {
//		panic(err)
//	}
//	code := m.Run()
//	os.Exit(code)
//}

func Test_redisPubSubConsumer_Consume(t *testing.T) {
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
				topics: "1",
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
				topics: "1",
				group:  "1",
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	mq := testPubSubMQ
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

func Test_redisPubSubProducer_Publish(t *testing.T) {
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
				msg: newRedisPubSubEvent("", []byte("t1m1")),
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
				msg: newRedisPubSubEvent("", []byte("t1m2")),
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
				p, err := testPubSubMQ.NewProducer(tt.args.ctx, tt.fields.topics...)
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
