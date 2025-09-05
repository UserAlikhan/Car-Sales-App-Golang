package utils

import (
	"car_sales/internal/configs"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DeleteFromS3(ctx context.Context, s3Conf *configs.S3Config, key string) error {
	_, err := s3Conf.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s3Conf.BucketName,
		Key:    &key,
	})
	if err != nil {
		return err
	}

	return nil
}

// this function generates presigned url that user can use on frontend
// to delete file by himself by simply clicking that url
func GetPresignedDeleteURL(ctx context.Context, s3Conf *configs.S3Config, key string, expire time.Duration) (string, error) {
	presignedClient := s3.NewPresignClient(s3Conf.Client)

	req, err := presignedClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s3Conf.BucketName,
		Key:    &key,
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", nil
	}

	return req.URL, nil
}
