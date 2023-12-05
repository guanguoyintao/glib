package uoss

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

type CosClient struct {
	*cos.Client
}

type CosOption func(c *CosClient)

func NewCosClient(cosUrl, cosSecretId, cosSecretKey string, options ...CosOption) *CosClient {
	// 构建cos客户端
	u, _ := url.Parse(cosUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cosSecretId,
			SecretKey: cosSecretKey,
		},
		Timeout: 5 * time.Minute,
	})
	cosClient := &CosClient{
		Client: c,
	}
	for _, opt := range options {
		opt(cosClient)
	}

	return cosClient
}

func (c *CosClient) UploadByFile(ctx context.Context, localFilePath, cosFilePath string, opts ...Option) (string, error) {
	if c == nil {
		return "", fmt.Errorf("no cos client init")
	}
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	contentType := "audio/wav"
	if len(opt.ContentType) > 0 {
		contentType = opt.ContentType
	}
	_, err := c.Object.PutFromFile(ctx, cosFilePath, localFilePath, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	})
	if err != nil {
		return "", err
	}
	fileUrl := c.Object.GetObjectURL(cosFilePath).String()

	return fileUrl, nil
}
