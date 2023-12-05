package ulocal

import (
	"fmt"
	"os"
	"testing"
)

func TestGetFd(t *testing.T) {
	type args struct {
		filePath string
		flag     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "not exit",
			args: args{
				filePath: "./tmp/test.txt",
				flag:     os.O_WRONLY | os.O_APPEND,
			},
			wantErr: false,
		},
		{
			name: "exit part",
			args: args{
				filePath: "./test/tmp/tmp/test.txt",
				flag:     os.O_WRONLY | os.O_APPEND,
			},
			wantErr: false,
		},
		{
			name: "exit",
			args: args{
				filePath: "./test/test.txt",
				flag:     os.O_WRONLY | os.O_APPEND,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFd(tt.args.filePath, tt.args.flag)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = got.Close()
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWriteHead(t *testing.T) {
	type args struct {
		head        []byte
		filePath    string
		newFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "equal",
			args: args{
				head:        []byte("hello word"),
				filePath:    "./test/test.txt",
				newFilePath: "./test/test.txt",
			},
			wantErr: false,
		},
		{
			name: "unequal",
			args: args{
				head:        []byte("hello word"),
				filePath:    "./test/test.txt",
				newFilePath: "./test/temp_test.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteHead(tt.args.head, tt.args.filePath, tt.args.newFilePath); (err != nil) != tt.wantErr {
				t.Errorf("WriteHead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveFileExtension(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "Case 1",
			args: "/a.t/b.t.t/c.t",
			want: "c",
		},
		{
			name: "Case 2",
			args: "/a/b/c",
			want: "c",
		},
		{
			name: "Case 3",
			args: "/a/b/c.t",
			want: "c",
		},
		{
			name: "Case 4",
			args: "c.t",
			want: "c",
		},
		{
			name: "Case 4",
			args: "c.t.t.t",
			want: "c.t.t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveFileExtension(tt.args)
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("RemoveFileExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
