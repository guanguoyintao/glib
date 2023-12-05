package uhttp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractFileNameFromURL(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "mp4",
			args: args{
				urlString: "https://umu-test-1303248253.cos.ap-beijing.myqcloud.com/resource/T/ab/x2xno/transcoding/433054406.webm.mp4",
			},
			want:    "433054406.webm.mp4",
			want1:   "/resource/T/ab/x2xno/transcoding",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ExtractFileNameFromURL(tt.args.urlString)
			if !tt.wantErr(t, err, fmt.Sprintf("ExtractFileNameFromURL(%v)", tt.args.urlString)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ExtractFileNameFromURL(%v)", tt.args.urlString)
			assert.Equalf(t, tt.want1, got1, "ExtractFileNameFromURL(%v)", tt.args.urlString)
		})
	}
}
