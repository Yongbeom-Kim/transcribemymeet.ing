package runpod_test

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/runpod"
)

var RUNPOD_WHISPER_URL = sync.OnceValue(func() string {
	url := os.Getenv("RUNPOD_WHISPER_URL")
	if url == "" {
		panic("RUNPOD_WHISPER_URL environment variable is not set")
	}
	return url
})

func TestRun(t *testing.T) {
	runResponse, err := runpod.Run(RUNPOD_WHISPER_URL(), runpod.RunRequest{
		Input: map[string]string{
			"audio": "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
			"model": "tiny",
		},
	})
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}

	switch runResponse.Status {
	case runpod.StatusQueue, runpod.StatusProgress, runpod.StatusComplete:
		t.Logf("Run response status: %v", runResponse.Status)
	default:
		t.Errorf("Invalid status: %v", runResponse.Status)
	}
	t.Logf("Run response: %v", runResponse)
}

func TestRunSync(t *testing.T) {
	runResponse, err := runpod.RunSync(RUNPOD_WHISPER_URL(), runpod.RunRequest{
		Input: map[string]string{
			"audio": "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
			"model": "tiny",
		},
	})
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	if runResponse.Status != runpod.StatusComplete {
		t.Errorf("Run resopnse status is not complete: %v", runResponse)
	}
	if runResponse.Output == nil {
		t.Errorf("Run response output is nil: %v", runResponse)
	}
	t.Logf("Run response: %v", runResponse)

	// Wait for job to complete
	for {
		time.Sleep(1 * time.Second)
		statusResponse, err := runpod.Status(RUNPOD_WHISPER_URL(), runResponse.JobId)
		if err != nil {
			t.Fatalf("Failed to job status: %v", err)
		}
		if statusResponse.Status == runpod.StatusComplete {
			break
		}
	}
}

func TestStatus(t *testing.T) {
	runResponse, err := runpod.Run(RUNPOD_WHISPER_URL(), runpod.RunRequest{
		Input: map[string]string{
			"audio": "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
			"model": "tiny",
		},
	})
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}

	statusResponse, err := runpod.Status(RUNPOD_WHISPER_URL(), runResponse.JobId)
	if err != nil {
		t.Errorf("Failed to get status: %v", err)
	}
	switch statusResponse.Status {
	case runpod.StatusQueue, runpod.StatusProgress, runpod.StatusComplete:
		t.Logf("Status response: %v", statusResponse)
	default:
		t.Errorf("Invalid status: %v", statusResponse.Status)
	}
}

func TestCancel(t *testing.T) {
	runResponse, err := runpod.Run(RUNPOD_WHISPER_URL(), runpod.RunRequest{
		Input: map[string]string{
			"audio": "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
			"model": "tiny",
		},
	})
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Run response: %v", runResponse)

	cancelResponse, err := runpod.Cancel(RUNPOD_WHISPER_URL(), runResponse.JobId)
	if err != nil {
		t.Errorf("Failed to cancel job: %v", err)
	}
	t.Logf("Cancel response: %v", cancelResponse)
}

func TestHealthCheck(t *testing.T) {
	healthCheckResponse, err := runpod.HealthCheck(RUNPOD_WHISPER_URL())
	if err != nil {
		t.Errorf("Failed to health check: %v", err)
	}
	t.Logf("Health check response: %v", healthCheckResponse)
}

func TestPurgeQueue(t *testing.T) {
	purgeResponse, err := runpod.PurgeQueue(RUNPOD_WHISPER_URL())
	if err != nil {
		t.Errorf("Failed to purge queue: %v", err)
	}
	t.Logf("Purge response: %v", purgeResponse)
}
