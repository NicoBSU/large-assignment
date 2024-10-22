package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"large-assignment/models"

	"github.com/minio/minio-go/v7"
)

type MinioService interface {
	CreateBucket(ctx context.Context, bucketName string) error
	PutObject(ctx context.Context, bucketName, objectName string, data []byte) error
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
}

type minioService struct {
	client *minio.Client
}

func NewMinioService(client *minio.Client) MinioService {
	return &minioService{client: client}
}

func (m *minioService) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket already exists: %w", err)
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

func (m *minioService) PutObject(ctx context.Context, bucketName, objectName string, data []byte) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to put object into bucket: %w", err)
	}
	return nil
}

func (m *minioService) GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	//interesting behavior, instead of returning an error in case if object was not found,
	//error message is stored in *minio.Object.
	if err != nil {
		return nil, fmt.Errorf("failed to get object from bucket: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, models.ERR_NOT_FOUND
		}
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return data, nil
}
