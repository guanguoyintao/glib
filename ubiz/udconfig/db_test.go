package udconfig

import (
	"context"
	"flag"
	"fmt"
	"git.umu.work/AI/uglib/ubiz"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/config"
	"git.umu.work/be/goframework/logger"
	"git.umu.work/be/goframework/store/gorm"
	"os"
	"path"
	"testing"
)

var dbConfig ubiz.UDConfig
var err error

type FeedbackMessageReviseTypeConfig struct {
	MessageReviseType        int    `json:"-"`
	MessageReviseTypeContent string `json:"message_revise_type_content"`
	ReviseTypeOptions        []int  `json:"revise_type_options"`
}

func messageReviseTypeDecoder(ctx context.Context, conf interface{}) (value interface{}, err error) {
	confString, ok := conf.(string)
	logger.GetLogger(ctx).Info(fmt.Sprintf("config is %s\n", confString))
	if !ok {
		return nil, uerrors.UErrorDynamicConfigTypeUnknown
	}
	var feedbackMessageReviseTypeConfig FeedbackMessageReviseTypeConfig
	err = ujson.Unmarshal([]byte(confString), &feedbackMessageReviseTypeConfig)
	if err != nil {
		return nil, err
	}

	return &feedbackMessageReviseTypeConfig, nil
}

func TestMain(m *testing.M) {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(currentPath)
	config.Init(path.Join(currentPath, "conf"))
	fmt.Printf("config init %+v", config.GetConfig())
	cache.Init(config.GetConfig())
	gorm.Init(config.GetConfig())
	dbConfig, err = NewDBDConfig(context.Background(), "ucs")
	if err != nil {
		panic(err)
	}
	flag.Parse()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestDBDConfig_Get(t *testing.T) {
	type args struct {
		ctx        context.Context
		key        string
		decodeFunc func(ctx context.Context, config interface{}) (value interface{}, err error)
	}
	tests := []struct {
		name      string
		args      args
		wantValue interface{}
		wantErr   bool
	}{
		{
			name: "feedback.message_revise_type.0",
			args: args{
				ctx:        context.Background(),
				key:        "feedback.message_revise_type.0",
				decodeFunc: messageReviseTypeDecoder,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dbConfig.RegisterDecoder(context.Background(), tt.args.key, tt.args.decodeFunc)
			if err != nil {
				t.Error(err)
			}
			gotValue, err := dbConfig.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotValue.(*FeedbackMessageReviseTypeConfig).MessageReviseTypeContent)
		})
	}
}
