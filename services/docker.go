package services

import (
	"context"
	"fmt"
	"large-assignment/config"
	"large-assignment/models"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerService interface {
	GetMinioInstances(ctx context.Context, conf config.Config) ([]types.Container, error)
	GetEnvVariables(ctx context.Context, container types.Container) (map[string]string, error)
	GetContainerIP(ctx context.Context, container types.Container, networkName string) (string, error)
	GetMinioConfigsFromInstances(ctx context.Context, conf config.Config) ([]models.MinioClientInfo, error)
}

type dockerService struct {
	dockerClient *client.Client
}

func NewDockerService(cli *client.Client) DockerService {
	return &dockerService{dockerClient: cli}
}

func InitDockerClient(ctx context.Context, conf config.Config) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(conf.DockerVersion))
	if err != nil {
		return nil, fmt.Errorf("failed to establish docker client: %w", err)
	}
	return cli, nil
}

func (s *dockerService) GetMinioInstances(ctx context.Context, conf config.Config) ([]types.Container, error) {
	containers, err := s.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get list of docker containers: %w", err)
	}

	var minioContainers []types.Container
	for _, container := range containers {
		if strings.Contains(container.Names[0], conf.MinioContainerNamePattern) {
			minioContainers = append(minioContainers, container)
		}
	}
	return minioContainers, nil
}

func (s *dockerService) GetEnvVariables(ctx context.Context, container types.Container) (map[string]string, error) {
	envVars := make(map[string]string)

	inspect, err := s.dockerClient.ContainerInspect(ctx, container.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	for _, env := range inspect.Config.Env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format - %s: %w", env, err)
		}
		envVars[parts[0]] = parts[1]
	}

	return envVars, nil
}

func (s *dockerService) GetContainerIP(ctx context.Context, container types.Container, networkName string) (string, error) {
	inspect, err := s.dockerClient.ContainerInspect(ctx, container.ID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	networkSettings := inspect.NetworkSettings.Networks[networkName]
	if networkSettings == nil {
		return "", fmt.Errorf("specified docker network not found: %w", err)
	}

	return networkSettings.IPAddress, nil
}

func (s *dockerService) GetMinioConfigsFromInstances(ctx context.Context, conf config.Config) ([]models.MinioClientInfo, error) {
	minioContainers, err := s.GetMinioInstances(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to get minio instances: %w", err)
	}

	var minioConfigs []models.MinioClientInfo
	for _, container := range minioContainers {
		envVars, err := s.GetEnvVariables(ctx, container)
		if err != nil {
			return nil, fmt.Errorf("failed to get env variables of container: %w", err)
		}

		ip, err := s.GetContainerIP(ctx, container, conf.DockerNetworkName)
		if err != nil {
			return nil, fmt.Errorf("failed to get container's exposed IP: %w", err)
		}

		accessKey, found := envVars[conf.UserEnvVarName]
		if !found {
			return nil, fmt.Errorf("container (%s) doesn't have such environment variable (%s): %w", container.ID, conf.UserEnvVarName, err)
		}

		secretKey, found := envVars[conf.PasswordEnvVarName]
		if !found {
			return nil, fmt.Errorf("container (%s) doesn't have such environment variable (%s): %w", container.ID, conf.PasswordEnvVarName, err)
		}

		minioConfig := models.MinioClientInfo{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Host:      ip,
			Port:      conf.MinioPort,
			Secure:    conf.MinioSecureConnection,
			Bucket:    conf.BucketName,
		}

		minioConfigs = append(minioConfigs, minioConfig)
	}

	return minioConfigs, nil
}
