package runpod

import (
	"log/slog"
	"os"
	"sync"
)

var RUNPOD_WHISPER_URL = sync.OnceValue(func() string {
	url := os.Getenv("RUNPOD_WHISPER_URL")
	if url == "" {
		panic("RUNPOD_WHISPER_URL environment variable is not set")
	}
	return url
})

func WhisperHealthCheck() {
	healthCheckResponse, err := HealthCheck(RUNPOD_WHISPER_URL())
	if err != nil {
		slog.Error("Error checking health", "error", err)
		return
	}

	slog.Info("Health check response", "response", healthCheckResponse)
}
