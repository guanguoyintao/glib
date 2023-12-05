package ujson

import (
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"io"
)

type jsonType string

const sonicType jsonType = "sonic"

var defaultJson uJsonAPI

func init() {
	// select json lib
	jsonT := sonicType
	switch jsonT {
	case sonicType:
		defaultJson = newSonicJsonAPI(sonic.Config{
			NoQuoteTextMarshaler: true,
		})
	}
}

// Marshal returns the JSON encoding bytes of v.
func Marshal(val interface{}) ([]byte, error) {
	return defaultJson.Marshal(val)
}

// MarshalString returns the JSON encoding string of v.
func MarshalString(val interface{}) (string, error) {
	return defaultJson.MarshalToString(val)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// NOTICE: This API copies given buffer by default,
// if you want to pass JSON more efficiently, use UnmarshalString instead.
func Unmarshal(buf []byte, val interface{}) error {
	return defaultJson.Unmarshal(buf, val)
}

// UnmarshalString is like Unmarshal, except buf is a string.
func UnmarshalString(buf string, val interface{}) error {
	return defaultJson.UnmarshalFromString(buf, val)
}

func DeleteJsonMapKey(in string, keys ...string) (string, error) {
	var m map[string]interface{}
	err := Unmarshal([]byte(in), &m)
	if err != nil {
		return "", errors.Wrapf(err, "json=%s", in)
	}
	for _, key := range keys {
		delete(m, key)
	}
	bs, err := Marshal(m)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func NewEncoder(writer io.Writer) Encoder {
	return defaultJson.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) Decoder {
	return defaultJson.NewDecoder(reader)
}
