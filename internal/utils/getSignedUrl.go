package utils

import (
	"car_sales/internal/configs"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetSignedUrl(ctx context.Context, s3Conf *configs.S3Config, bucketName string, key string, expire time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(s3Conf.Client)
	// images by default are private and could be only accessed by designated user
	// however, we can give an temporary access using PresignGetObject
	req, err := presignedClient.PresignGetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			// ResponseContentDisposition: aws.String("inline"),
		},
		s3.WithPresignExpires(expire),
	)
	if err != nil {
		return "", err
	}

	return req.URL, err
}
