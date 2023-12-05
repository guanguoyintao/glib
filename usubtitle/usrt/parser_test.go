package usrt

import (
	"context"
	"testing"
)

func TestWriteFile(t *testing.T) {
	type args struct {
		ctx        context.Context
		srtPath    string
		outSrtPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				ctx:        nil,
				srtPath:    "/Users/guanguoyintao/work/golang/uglib/usubtitle/usrt/3962881944.srt",
				outSrtPath: "/Users/guanguoyintao/work/golang/uglib/usubtitle/usrt/3962881944.output.srt",
			},
		},
		{
			name: "2",
			args: args{
				ctx:        nil,
				srtPath:    "/Users/guanguoyintao/work/golang/uglib/usubtitle/usrt/4074389319.srt",
				outSrtPath: "/Users/guanguoyintao/work/golang/uglib/usubtitle/usrt/4074389319.output.srt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := ParserFromFile(tt.args.ctx, tt.args.srtPath)
			if err != nil {
				t.Error(err.Error())
				return
			}
			err = WriteFile(tt.args.ctx, tt.args.outSrtPath, items)
			if err != nil {
				t.Error(err.Error())
				return
			}
		})
	}
}
