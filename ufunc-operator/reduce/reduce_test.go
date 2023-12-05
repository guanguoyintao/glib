package ureduce

import (
	"context"
	"reflect"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestReduce(t *testing.T) {
	type args struct {
		ctx          context.Context
		slice        interface{}
		fn           func(ctx context.Context, acc, curr interface{}) (interface{}, error)
		initialValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Reduce with an array of integers",
			args: args{
				ctx:   context.TODO(),
				slice: []int{1, 2, 3, 4, 5},
				fn: func(ctx context.Context, acc, curr interface{}) (interface{}, error) {
					return acc.(int) + curr.(int), nil
				},
				initialValue: 0,
			},
			want:    15,
			wantErr: false,
		},
		{
			name: "Reduce with an array of strings",
			args: args{
				ctx:   context.TODO(),
				slice: []string{"Hello", " ", "World", "!"},
				fn: func(ctx context.Context, acc, curr interface{}) (interface{}, error) {
					return acc.(string) + curr.(string), nil
				},
				initialValue: "",
			},
			want:    "Hello World!",
			wantErr: false,
		},
		{
			name: "Reduce with an array of pointers",
			args: args{
				ctx:   context.TODO(),
				slice: []*int{new(int), new(int), new(int)},
				fn: func(ctx context.Context, acc, curr interface{}) (interface{}, error) {
					return acc.(int) + *curr.(*int), nil
				},
				initialValue: 0,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Reduce with an array of custom structs",
			args: args{
				ctx: context.TODO(),
				slice: []Person{
					{Name: "Alice", Age: 25},
					{Name: "Bob", Age: 30},
					{Name: "Charlie", Age: 35},
				},
				fn: func(ctx context.Context, acc, curr interface{}) (interface{}, error) {
					person := curr.(Person)
					return acc.(int) + person.Age, nil
				},
				initialValue: 0,
			},
			want:    90,
			wantErr: false,
		},
		{
			name: "Reduce with an empty slice",
			args: args{
				ctx:   context.TODO(),
				slice: []int{},
				fn: func(ctx context.Context, acc, curr interface{}) (interface{}, error) {
					return acc.(int) + curr.(int), nil
				},
				initialValue: 0,
			},
			want:    0,
			wantErr: false,
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Reduce(tt.args.ctx, tt.args.slice, tt.args.fn, tt.args.initialValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reduce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() got = %v, want %v", got, tt.want)
			}
		})
	}
}
