package filter

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/utypes/uslice"
	"git.umu.work/be/goframework/logger"
)

func Filter(ctx context.Context, slice interface{}, predicate func(ctx context.Context, item interface{}) (bool, error)) (interface{}, error) {
	if slice == nil {
		return nil, nil
	}
	s, ok := uslice.Convert2Slice(slice)
	if !ok {
		return nil, uerrors.UErrorFuncOperatorTypeAssertion
	}
	var filteredSlice []interface{}
	for _, v := range s {
		ok, err := predicate(ctx, v)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
		if ok {
			filteredSlice = append(filteredSlice, v)
		}
	}

	return filteredSlice, nil
}
