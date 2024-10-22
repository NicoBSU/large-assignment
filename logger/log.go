package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.Logger

func NewLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	Log = logger
}
