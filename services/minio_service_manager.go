package services

import (
	"context"
	"fmt"
	"large-assignment/models"
	"large-assignment/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioServiceManager interface {
	GetMinioService(id string) MinioService
}

type minioServiceManager struct {
	services []MinioService
}

func NewMinioServiceManager(ctx context.Context, configs []models.MinioClientInfo) (MinioServiceManager, error) {
	services := make([]MinioService, len(configs))

	for i, config := range configs {
		endpoint := fmt.Sprintf("%s:%s", config.Host, config.Port)
		client, err := minio.New(endpoint,
			&minio.Options{
				Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
				Secure: config.Secure,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to establish minio client: %w", err)
		}
		services[i] = NewMinioService(client)
		err = services[i].CreateBucket(ctx, config.Bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket for specified minio instance (%s): %w", endpoint, err)
		}
	}

	return &minioServiceManager{services: services}, nil
}

func (m *minioServiceManager) GetMinioService(id string) MinioService {
	instanceIndex := util.GetMinioInstanceFromId(id, len(m.services))
	service := m.services[instanceIndex]
	return service
}
