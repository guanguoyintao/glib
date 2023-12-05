package urand

import (
	"fmt"
	"testing"
)

func TestGenStr(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "100",
			args: args{
				n: 1,
			},
		},
		{
			name: "100",
			args: args{
				n: 2,
			},
		},
		{
			name: "100",
			args: args{
				n: 10,
			},
		},
		{
			name: "100",
			args: args{
				n: 10,
			},
		},
		{
			name: "100",
			args: args{
				n: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenStr(tt.args.n)
			fmt.Println(got)
		})
	}
}
