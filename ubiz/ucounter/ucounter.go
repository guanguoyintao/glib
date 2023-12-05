package ucounter

import (
	"context"
	"git.umu.work/AI/uglib/ubiz"
)

func NewCounter(ctx context.Context, namespace CounterNameSpaceType) (ubiz.UCounter, error) {
	return NewDBCounter(ctx, namespace)
}
