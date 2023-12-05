package ujson

import (
	"encoding/json"
	"io"
)

type jsonApi struct{}

func newJsonAPI() uJsonAPI {
	return &jsonApi{}
}

func (j jsonApi) MarshalToString(v interface{}) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (j jsonApi) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonApi) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (j jsonApi) UnmarshalFromString(str string, v interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (j jsonApi) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j jsonApi) NewEncoder(writer io.Writer) Encoder {
	//TODO implement me
	panic("implement me")
}

func (j jsonApi) NewDecoder(reader io.Reader) Decoder {
	//TODO implement me
	panic("implement me")
}

func (j jsonApi) Valid(data []byte) bool {
	//TODO implement me
	panic("implement me")
}
