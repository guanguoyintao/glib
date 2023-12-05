package uhash

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashMD532(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashMD532(tt.args.s)
			fmt.Println(got)
			assert.Equal(t, 32, len(got))
		})
	}
}

func TestHashMurmurHash340(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				s: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashMurmurHash340(tt.args.s)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(got)
			assert.Equal(t, 40, len(got))
		})
	}
}

func TestCalcFileSHA256(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ExistingFile",
			args: args{
				filePath: "existing.txt", // 指定一个已存在的文件路径
			},
			want:    "3bace75732dae79185fed42047d71dfb32a4cb90605d87d4cf03e2d14b09f3d7",
			wantErr: assert.NoError,
		},
		{
			name: "EmptyFile",
			args: args{
				filePath: "empty.txt", // 指定一个空文件路径
			},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr: assert.NoError,
		},
		{
			name: "NonExistingFile",
			args: args{
				filePath: "non_existing.txt", // 指定一个不存在的文件路径
			},
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcFileSHA256(context.Background(), tt.args.filePath)
			if !tt.wantErr(t, err, fmt.Sprintf("CalculateFileSHA256(%v)", tt.args.filePath)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CalculateFileSHA256(%v)", tt.args.filePath)
			if err == nil {
				assert.Equalf(t, 64, len(got), "CalculateFileSHA256(%v)", tt.args.filePath)
			}
		})
	}
}
