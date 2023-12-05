package umap

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/utypes/uslice"
	"git.umu.work/be/goframework/logger"
)

func FlatMap(ctx context.Context, slice interface{}, handler func(ctx context.Context, item interface{}) (interface{}, error)) ([]interface{}, error) {
	if slice == nil {
		return nil, nil
	}
	s, ok := uslice.FlattenNestedSlice(slice)
	if !ok {
		return nil, uerrors.UErrorFuncOperatorTypeAssertion
	}
	var result []interface{}
	for _, v := range s {
		mapped, err := handler(ctx, v)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
		switch mv := mapped.(type) {
		case []interface{}:
			result = append(result, mv...)
		default:
			result = append(result, mapped)
		}
	}
	return result, nil
}
