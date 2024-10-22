package main

import (
	"context"
	"errors"
	"fmt"
	"large-assignment/config"
	"large-assignment/handlers"
	"large-assignment/logger"
	"large-assignment/router"
	"large-assignment/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.NewLogger()
	logger.Log.Info("Config loaded, logger initialized")

	dockerClient, err := services.InitDockerClient(ctx, *conf)
	if err != nil {
		logger.Log.Fatal("Failed to init docker client: ", zap.Error(err))
	}

	dockerService := services.NewDockerService(dockerClient)
	minioConfigs, err := dockerService.GetMinioConfigsFromInstances(ctx, *conf)
	if err != nil {
		logger.Log.Fatal("Failed to retrieve minio configs from docker containers: ", zap.Error(err))
	}

	clientManager, err := services.NewMinioServiceManager(ctx, minioConfigs)
	if err != nil {
		logger.Log.Fatal("Failed to init minio services manager: ", zap.Error(err))
	}

	handler := handlers.NewHandler(clientManager, conf.BucketName)

	r := router.InitRoutes(*handler)

	server := http.Server{
		Addr:    conf.AppHost + conf.AppPort,
		Handler: r,
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		<-ctx.Done()
		log.Println("Shut down gracefully...")

		defer func() {
			stop()
			dockerClient.Close()
			logger.Log.Sync()

			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Fatalf("Failed to shutdown http server: %s", err)
			}

			log.Println("Graceful shutdown completed")
			close(errC)
		}()

	}()

	go func() {
		logger.Log.Info(fmt.Sprintf("Listening on %s%s", conf.AppHost, conf.AppPort))
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	err = <-errC
	if err != nil {
		log.Fatalf("Server exited with error: %s", err)
	}
}
