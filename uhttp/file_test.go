package uhttp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetFileInfo(t *testing.T) {
	type args struct {
		fileURL string
	}
	tests := []struct {
		name    string
		args    args
		want    *FileInfo
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "cos",
			args: args{
				fileURL: "http://umu-cn-1303248253.cos.ap-beijing.myqcloud.com/4074389319.zh.srt",
			},
			want: &FileInfo{
				Size:         954,
				Type:         FileTypeText,
				Name:         "4074389319.zh.srt",
				ModifiedTime: time.Date(2023, time.October, 25, 6, 48, 33, 0, time.FixedZone("GMT", 0)),
			},
			wantErr: assert.NoError,
		},
		{
			name: "umu-cn.umucdn.cn",
			args: args{
				fileURL: "http://umu-cn.umucdn.cn/4074389319.zh.srt",
			},
			want: &FileInfo{
				Size:         954,
				Type:         FileTypeText,
				Name:         "4074389319.zh.srt",
				ModifiedTime: time.Date(2023, time.October, 25, 6, 48, 33, 0, time.FixedZone("GMT", 0)),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileInfo, err := GetFileInfo(tt.args.fileURL)
			if !tt.wantErr(t, err, fmt.Sprintf("GetFileInfo(%v)", tt.args.fileURL)) {
				return
			}
			assert.Equalf(t, tt.want, fileInfo, "GetFileInfo(%v)", tt.args.fileURL)
		})
	}
}
