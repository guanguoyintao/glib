package udconfig

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/ujson"
)

func JsonDecoder(ctx context.Context, conf interface{}) (value interface{}, err error) {
	var buf []byte
	switch conf.(type) {
	case string:
		buf = []byte(conf.(string))
	case []byte:
		buf = conf.([]byte)
	default:
		return nil, uerrors.UErrorDynamicConfigTypeUnknown
	}
	err = ujson.Unmarshal(buf, value)
	if err != nil {
		return nil, err
	}

	return value, err
}
