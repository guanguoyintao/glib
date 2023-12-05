package umap

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/utypes/uslice"
	"git.umu.work/be/goframework/logger"
)

func Map(ctx context.Context, slice interface{}, handler func(ctx context.Context, item interface{}) (interface{}, error)) (interface{}, error) {
	if slice == nil {
		return nil, nil
	}
	s, ok := uslice.Convert2Slice(slice)
	if !ok {
		return nil, uerrors.UErrorFuncOperatorTypeAssertion
	}
	result := make([]interface{}, len(s))
	for i, v := range s {
		var err error
		result[i], err = handler(ctx, v)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
	}
	return result, nil
}
