package uffmpeg

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetStreamingMediaInfo(t *testing.T) {
	type args struct {
		ctx       context.Context
		url       string
		mediaType StreamingMediaType
	}
	tests := []struct {
		name    string
		args    args
		want    *StreamingMediaInfo
		wantErr bool
	}{
		{
			name: "mp4",
			args: args{
				ctx:       context.Background(),
				url:       "https://umu-cn.umucdn.cn/resource/1i/OVO5/tpMUq/1010786800.mp4",
				mediaType: VideoStreamingMediaType,
			},
			want: &StreamingMediaInfo{
				CodecName: "h264",
				Duration:  time.Duration(1214.186625 * float64(time.Second)),
			},
			wantErr: false,
		},
		{
			name: "without image mp4",
			args: args{
				ctx:       context.Background(),
				url:       "http://umu-test.umucdn.cn/ushow/nuoheyi7_test.mp4",
				mediaType: VideoStreamingMediaType,
			},
			want: &StreamingMediaInfo{
				CodecName: "h264",
				Duration:  time.Duration(1214.186625 * float64(time.Second)),
			},
			wantErr: false,
		},
		{
			name: "local audio mp3",
			args: args{
				ctx:       context.Background(),
				url:       "/Users/guanguoyintao/Downloads/1.mp3",
				mediaType: AudioStreamingMediaType,
			},
			want: &StreamingMediaInfo{
				CodecName: "h264",
				Duration:  time.Duration(1214.186625 * float64(time.Second)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell := os.Getenv("SHELL")
			fmt.Println("$SHELL 环境变量的值是:", shell)
			p := os.Getenv("PATH")
			fmt.Println("$PATH 环境变量的值是:", p)
			got, err := GetStreamingMediaInfo(tt.args.ctx, tt.args.url, tt.args.mediaType)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteStreamingMediaInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteStreamingMediaInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
