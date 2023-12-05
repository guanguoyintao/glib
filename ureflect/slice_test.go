package ureflect

import (
	"reflect"
	"testing"
)

func TestConvertInterfaceToSlice(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  []interface{}
		want1 bool
	}{
		{
			name: "int slice",
			args: args{
				input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want:  []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9},
			want1: true,
		},
		{
			name: "string slice",
			args: args{
				input: []string{"a", "b", "c"},
			},
			want:  []interface{}{"a", "b", "c"},
			want1: true,
		},
		{
			name: "nested slice",
			args: args{
				input: []interface{}{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []interface{}{"a", "b", "c"}},
			},
			want:  []interface{}{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []interface{}{"a", "b", "c"}},
			want1: true,
		},
		{
			name: "struct",
			args: args{
				input: struct{}{},
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ConvertInterfaceToSlice(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertInterfaceToSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ConvertInterfaceToSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
