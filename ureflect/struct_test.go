package ureflect

import (
	"git.umu.work/AI/uglib/utest"
	"reflect"
	"testing"
)

func TestGetPackageName(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "struct",
			args: args{
				obj: utest.B{},
			},
			want: "git.umu.work/AI/uglib/utest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPackageName(tt.args.obj)
			if got != tt.want {
				t.Errorf("GetPackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStruct2Map(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "struct",
			args: args{
				obj: utest.B{},
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Struct2Map(tt.args.obj)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Struct2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMethodNames(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "struct",
			args: args{
				obj: &utest.B{},
			},
			want: []string{"GetPreheat", "ResetStartTime"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMethodNames(tt.args.obj)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMethodNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStructName(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "struct",
			args: args{
				obj: utest.B{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStructName(tt.args.obj)
			if got != tt.want {
				t.Errorf("GetStructName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMethodMap(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]func()
	}{
		{
			name: "struct",
			args: args{
				obj: &utest.B{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMethodMap(tt.args.obj)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMethodMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
