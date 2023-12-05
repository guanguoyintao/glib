package uoss

import (
	gstorage "cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"io"
	"os"
)

type GCSClient struct {
	*gstorage.Client
	bucket string
}

func NewGCSClient(ctx context.Context, credentials, bucket string) (*GCSClient, error) {
	cli, err := gstorage.NewClient(ctx, option.WithCredentialsJSON([]byte(credentials)))
	if err != nil {
		return nil, err
	}

	return &GCSClient{
		Client: cli,
		bucket: bucket,
	}, nil
}

func (g *GCSClient) UploadByFile(ctx context.Context, localFilePath, gcsFilePath string, opts ...Option) (string, error) {
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	if g == nil {
		return "", fmt.Errorf("no gcs client init")
	}
	fd, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	// Upload an object with storage.Writer.
	wc := g.Client.Bucket(g.bucket).Object(gcsFilePath).NewWriter(ctx)
	if _, err = io.Copy(wc, fd); err != nil {
		return "", errors.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", errors.Errorf("Writer.Close: %v", err)
	}
	return fmt.Sprintf("gs://%s/%s", g.bucket, gcsFilePath), nil
}
