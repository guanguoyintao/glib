package umap

import (
	"context"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		ctx     context.Context
		slice   interface{}
		handler Handler
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "string map",
			args: args{
				ctx:   context.Background(),
				slice: []string{"a", "b", "c"},
				handler: func(ctx context.Context, item interface{}) (interface{}, error) {
					i := item.(string)

					return i + "_map", nil
				},
			},
			want:    []string{"a_map", "b_map", "c_map"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map(tt.args.ctx, tt.args.slice, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.IsEqual(got, tt.want)
		})
	}
}
