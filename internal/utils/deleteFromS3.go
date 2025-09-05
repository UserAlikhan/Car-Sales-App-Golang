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

	waiter := s3.NewObjectExistsWaiter(s3Conf.Client)
	return waiter.Wait(ctx, &s3.HeadObjectInput{
		Bucket: &s3Conf.BucketName,
		Key:    &key,
	}, 30*time.Second)
}
