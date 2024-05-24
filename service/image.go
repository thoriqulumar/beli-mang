package service

import (
	"beli-mang/config"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

type ImageService interface {
	UploadImage(file *multipart.FileHeader) <-chan string
}

type imageService struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewImageService(cfg *config.Config, logger *zap.Logger) ImageService {
	return &imageService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *imageService) UploadImage(file *multipart.FileHeader) <-chan string {
	fileURLChan := make(chan string, 1)

	go func() {
		defer close(fileURLChan)

		src, err := file.Open()
		if err != nil {
			s.logger.Error("Failed to open file for upload:", zap.Error(err))
			fileURLChan <- ""
			return
		}
		defer src.Close()

		uuid := uuid.New().String()
		fileName := uuid + ".jpeg"

		url, err := uploadToS3(src, fileName, s.cfg)
		if err != nil {
			s.logger.Error("Failed to upload file to S3:", zap.Error(err))
			fileURLChan <- ""
			return
		}

		fileURLChan <- url
	}()

	return fileURLChan
}

func uploadToS3(file io.Reader, filename string, cfg *config.Config) (string, error) {
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.S3Region),
		Credentials: credentials.NewStaticCredentials(cfg.S3AcessKey, cfg.S3Secret, ""),
	})
	if err != nil {
		return "", errors.New("failed to create AWS session")
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Specify bucket name and object key
	bucketName := cfg.S3Bucket
	objectKey := filename

	// Upload file to S3
	_, err = svc.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("awss3." + objectKey),
		ACL:    aws.String("public-read"),
		Body:   aws.ReadSeekCloser(file),
	})
	if err != nil {
		return "", errors.New("failed to upload file to S3")
	}

	// Generate S3 object URL
	objectURL := fmt.Sprintf("https://awss3.%s", objectKey)

	return objectURL, nil
}
