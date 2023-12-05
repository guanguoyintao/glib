package ujson

import (
	"fmt"
	"testing"
)

type testJsonType struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Request struct {
		GroupId string `json:"group_id"`
	} `json:"request"`
}

func TestMarshal(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				val: &testJsonType{
					Service: "1",
					Method:  "2",
					Request: struct {
						GroupId string `json:"group_id"`
					}{
						GroupId: "3",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(string(got))
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		buf []byte
		val interface{}
	}

	var data testJsonType

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				buf: []byte(`{"service":"1","method":"2","request":{"group_id":"3"}}`),
				val: &data,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.args.buf, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("%+v\n", data)
		})
	}
}
