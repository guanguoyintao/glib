package umetadata

import (
	"context"
	"git.umu.work/be/goframework/metadata"
	gmetadata "github.com/micro/go-micro/v2/metadata"
)

func ToMap(ctx context.Context) map[string]string {
	if ctx == nil {
		ctx = context.Background()
	}
	md, ok := gmetadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	mmd, ok := metadata.FromContext(ctx)
	if !ok {
		return make(map[string]string)
	}
	for k, v := range mmd {
		md[k] = v
	}

	return md
}

func Merge(ctx context.Context, md map[string]string) context.Context {
	if md == nil {
		return ctx
	}
	ctx = metadata.MergeContext(ctx, md, true)
	return ctx
}

func ClientCtx(ctx context.Context) context.Context {
	header := ToMap(ctx)
	ctx = gmetadata.MergeContext(ctx, header, true)
	ctx = metadata.MergeContext(ctx, header, true)
	return ctx
}
