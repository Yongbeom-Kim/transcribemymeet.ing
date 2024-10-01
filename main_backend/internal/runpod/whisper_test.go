package runpod_test

import (
	"testing"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/runpod"
)

var input = runpod.WhisperModelInput{
	Audio: "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
	Model: "tiny",
}

func TestWhisperRun(t *testing.T) {
	response, err := runpod.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperRunSync(t *testing.T) {
	response, err := runpod.WhisperRunSync(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperStatus(t *testing.T) {
	response, err := runpod.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	status, err := runpod.WhisperStatus(response.JobId)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	switch status.Status {
	case runpod.StatusQueue, runpod.StatusProgress, runpod.StatusComplete:
		t.Logf("Status: %v", status)
	default:
		t.Errorf("Invalid status: %v", status)
	}
}

func TestWhisperCancel(t *testing.T) {
	response, err := runpod.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	cancel, err := runpod.WhisperCancel(response.JobId)
	if err != nil {
		t.Fatalf("Failed to cancel job: %v", err)
	}
	t.Logf("Cancel: %v", cancel)
}

func TestWhisperHealthCheck(t *testing.T) {
	response, err := runpod.WhisperHealthCheck()
	if err != nil {
		t.Fatalf("Failed to health check: %v", err)
	}
	t.Logf("Health check: %v", response)
}

func TestWhisperPurgeQueue(t *testing.T) {
	response, err := runpod.WhisperPurgeQueue()
	if err != nil {
		t.Fatalf("Failed to purge queue: %v", err)
	}
	t.Logf("Purge queue: %v", response)
}
