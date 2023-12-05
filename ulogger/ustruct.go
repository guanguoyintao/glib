package ulogger

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/be/goframework/logger"
)

func StructString(ctx context.Context, s interface{}) string {
	bs, err := ujson.Marshal(s)
	if err != nil {
		logger.GetLogger(ctx).Error(fmt.Sprintf("struct %+v to json error", s))

		return ""
	}

	return string(bs)
}
