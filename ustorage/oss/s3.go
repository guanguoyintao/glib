package uoss

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type S3Client struct {
	*s3.Client
	s3Bucket string
}

func NewS3Client(s3Region, s3Bucket, s3AccessKey, s3SecretKey string) *S3Client {
	// 构建s3 client
	cfg := aws.Config{
		Region:      s3Region,
		Credentials: credentials.NewStaticCredentialsProvider(s3AccessKey, s3SecretKey, ""),
	}
	client := &S3Client{
		Client:   s3.NewFromConfig(cfg),
		s3Bucket: s3Bucket,
	}

	return client
}

func (s *S3Client) UploadByFile(ctx context.Context, localFilePath, s3FilePath string, opts ...Option) (string, error) {
	if s == nil {
		return "", fmt.Errorf("no s3 client init")
	}
	fd, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	opt := &Options{}
	for _, o := range opts {
		o(opt)
	}
	contentType := "audio/wav"
	if len(opt.ContentType) > 0 {
		contentType = opt.ContentType
	}
	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Bucket),
		Key:         aws.String(s3FilePath),
		Body:        fd,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}
	locationOutput, err := s.Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(s.s3Bucket),
	})
	if err != nil {
		return "", err
	}
	region := string(locationOutput.LocationConstraint)
	fileUrl := fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", region, s.s3Bucket, s3FilePath)

	return fileUrl, nil
}
