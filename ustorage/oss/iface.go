package uoss

import "context"

type ObjectStorage interface {
	UploadByFile(ctx context.Context, localFilePath, uosFilePath string, opts ...Option) (string, error)
}
