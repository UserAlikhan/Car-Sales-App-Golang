package configs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	BucketName string
	Region     string
	Client     *s3.Client
	Uploader   *manager.Uploader
}

func NewS3Config(ctx context.Context, bucketName string, region string) (*S3Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Config{
		BucketName: bucketName,
		Region:     region,
		Client:     client,
		Uploader:   uploader,
	}, nil
}
