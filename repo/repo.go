package repo

import (
	"context"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3Serv "github.com/aws/aws-sdk-go/service/s3"
)

const (
	s3PublicACL = "public-read"
)

func NewS3Repository(ctx context.Context, region, url, accessKey, secretKey string) (S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(url),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return &s3{
		url: url,
		svc: s3Serv.New(sess),
	}, nil

}

type S3 interface {
	Upload(ctx context.Context, file io.ReadSeeker, fileName, bucket string) (string, error)
	Delete(ctx context.Context, fileUrl, bucket string) error
	GetURL(ctx context.Context, bucketName, fileName string) (string, error)
}

type s3 struct {
	url string
	svc *s3Serv.S3
}

func (s *s3) Upload(ctx context.Context, file io.ReadSeeker, fileName, bucket string) (string, error) {
	if _, err := s.svc.PutObjectWithContext(ctx, &s3Serv.PutObjectInput{
		Body:               file,
		Bucket:             aws.String(bucket),
		Key:                aws.String(fileName),
		ACL:                aws.String(s3PublicACL),
		ContentDisposition: aws.String("inline"),
		ContentType:        aws.String("audio/ogg"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", s.url, bucket, fileName), nil
}

func (s *s3) GetURL(ctx context.Context, bucketName, fileName string) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3Serv.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	url, err := req.Presign(7 * 24 * time.Hour)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *s3) Delete(ctx context.Context, fileUrl, bucket string) error {
	_, err := s.svc.DeleteObjectWithContext(ctx, &s3Serv.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path.Base(fileUrl)),
	})
	return err
}
