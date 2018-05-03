package main

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	Bucket     string
	S3Client   *s3.S3
	S3Uploader *s3manager.Uploader
}

func NewS3Storage(bucket, accessKey, secretKey string) *S3Storage {
	sess := session.Must(session.NewSession())
	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	awsOptions := &aws.Config{
		Credentials: creds,
		Region:      aws.String(endpoints.EuCentral1RegionID),
	}
	svc := s3.New(sess, awsOptions)

	uploader := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		u.PartSize = 8 * 1024 * 1024
	})

	return &S3Storage{
		Bucket:     bucket,
		S3Client:   svc,
		S3Uploader: uploader,
	}
}

func (s *S3Storage) ReadFile(path string) (io.ReadCloser, error) {
	ctx := context.Background()
	result, err := s.S3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *S3Storage) WriteFile(path string, file io.ReadCloser) error {
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
		Body:   file,
	}

	_, err := s.S3Uploader.Upload(upParams)
	return err
}
