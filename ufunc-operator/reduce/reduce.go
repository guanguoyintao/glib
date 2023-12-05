package ureduce

import (
	"context"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/utypes/uslice"
)

func Reduce(ctx context.Context, slice interface{}, fn func(ctx context.Context, acc, curr interface{}) (interface{}, error), initialValue interface{}) (interface{}, error) {
	if slice == nil {
		return initialValue, nil
	}
	s, ok := uslice.Convert2Slice(slice)
	if !ok {
		return nil, uerrors.UErrorFuncOperatorTypeAssertion
	}
	acc := initialValue
	for _, v := range s {
		var err error
		acc, err = fn(ctx, acc, v)
		if err != nil {
			return nil, err
		}
	}

	return acc, nil
}
