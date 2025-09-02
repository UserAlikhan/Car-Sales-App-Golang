package utils

import (
	"car_sales/internal/configs"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadToS3(
	ctx *gin.Context, s3Conf *configs.S3Config, key *string,
	file multipart.File, contentType string,
) (*manager.UploadOutput, error) {
	output, err := s3Conf.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      &s3Conf.BucketName,
		Key:         key,
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}
