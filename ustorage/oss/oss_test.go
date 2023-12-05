package uoss

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCovertUrl2Outer(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "cos:umu-cn-1303248253.cos.ap-beijing.myqcloud.com",
			args: args{
				urlString: "http://umu-cn-1303248253.cos.ap-beijing.myqcloud.com/4074389319.zh.srt",
			},
			want:    "http://umu-cn.umucdn.cn/4074389319.zh.srt",
			wantErr: assert.NoError,
		},
		{
			name: "s3-co",
			args: args{
				urlString: "https://s3.ap-northeast-1.amazonaws.com/umu.co/videoweike/teacher/weike/20ebe/transcoding/1700192712.1681.19144.202260481.mp4.mp4",
			},
			want:    "https://co.umustatic.com/videoweike/teacher/weike/20ebe/transcoding/1700192712.1681.19144.202260481.mp4.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "s3-tw",
			args: args{
				urlString: "https://s3.ap-northeast-1.amazonaws.com/umu.tw/videoweike/teacher/weike/7ebea/transcoding/1699860032.965.44429.200132985.mp4.mp4",
			},
			want:    "https://tw.umustatic.com/videoweike/teacher/weike/7ebea/transcoding/1699860032.965.44429.200132985.mp4.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "s3-io",
			args: args{
				urlString: "https://s3.eu-west-1.amazonaws.com/umu.io/videoweike/teacher/weike/IZu1775/transcoding/1699929844.1392.92528.200045584.mov.mp4",
			},
			want:    "https://resource.umu.io/videoweike/teacher/weike/IZu1775/transcoding/1699929844.1392.92528.200045584.mov.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "s3-com",
			args: args{
				urlString: "https://s3.us-west-2.amazonaws.com/umu.com/videoweike/teacher/weike/EVc9e4/transcoding/1699883327.7847.34342.200135144.mp4.mp4",
			},
			want:    "https://com.umustatic.com/videoweike/teacher/weike/EVc9e4/transcoding/1699883327.7847.34342.200135144.mp4.mp4",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CovertUrl2Outer(tt.args.urlString)
			if !tt.wantErr(t, err, fmt.Sprintf("CovertUrl2Outer(%v)", tt.args.urlString)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CovertUrl2Outer(%v)", tt.args.urlString)
		})
	}
}

func TestCovertUrl2Inner(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "cn",
			args: args{
				urlString: "http://umu-cn.umucdn.cn/4074389319.zh.srt",
			},
			want:    "http://umu-cn-1303248253.cos.ap-beijing.myqcloud.com/4074389319.zh.srt",
			wantErr: assert.NoError,
		},
		{
			name: "co",
			args: args{
				urlString: "https://co.umustatic.com/videoweike/teacher/weike/20ebe/transcoding/1700192712.1681.19144.202260481.mp4.mp4",
			},
			want:    "https://s3.ap-northeast-1.amazonaws.com/umu.co/videoweike/teacher/weike/20ebe/transcoding/1700192712.1681.19144.202260481.mp4.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "tw",
			args: args{
				urlString: "https://tw.umustatic.com/videoweike/teacher/weike/7ebea/transcoding/1699860032.965.44429.200132985.mp4.mp4",
			},
			want:    "https://s3.ap-northeast-1.amazonaws.com/umu.tw/videoweike/teacher/weike/7ebea/transcoding/1699860032.965.44429.200132985.mp4.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "io",
			args: args{
				urlString: "https://resource.umu.io/videoweike/teacher/weike/IZu1775/transcoding/1699929844.1392.92528.200045584.mov.mp4",
			},
			want:    "https://s3.eu-west-1.amazonaws.com/umu.io/videoweike/teacher/weike/IZu1775/transcoding/1699929844.1392.92528.200045584.mov.mp4",
			wantErr: assert.NoError,
		},
		{
			name: "com",
			args: args{
				urlString: "https://com.umustatic.com/videoweike/teacher/weike/EVc9e4/transcoding/1699883327.7847.34342.200135144.mp4.mp4",
			},
			want:    "https://s3.us-west-2.amazonaws.com/umu.com/videoweike/teacher/weike/EVc9e4/transcoding/1699883327.7847.34342.200135144.mp4.mp4",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CovertUrl2Inner(tt.args.urlString)
			if !tt.wantErr(t, err, fmt.Sprintf("CovertUrl2Inner(%v)", tt.args.urlString)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CovertUrl2Inner(%v)", tt.args.urlString)
		})
	}
}
