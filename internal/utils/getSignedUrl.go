package utils

import (
	"car_sales/internal/configs"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetSignedUrl(ctx context.Context, s3Conf *configs.S3Config, bucketName string, key string, expire time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(s3Conf.Client)

	req, err := presignedClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", err
	}

	return req.URL, err
}
