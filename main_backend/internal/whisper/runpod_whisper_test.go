package whisper_test

import (
	"testing"
	"time"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/whisper"
)

var input = whisper.NewWhisperInput(
	"https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
	whisper.WithModel(whisper.WhisperModelTiny),
)

func TestWhisperRun(t *testing.T) {
	response, err := whisper.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperRunSync(t *testing.T) {
	response, err := whisper.WhisperRunSync(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperStatus(t *testing.T) {
	response, err := whisper.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	status, err := whisper.WhisperStatus(response.JobId)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	switch status.Status {
	case whisper.StatusQueue, whisper.StatusProgress, whisper.StatusComplete:
		t.Logf("Status: %v", status)
	default:
		t.Errorf("Invalid status: %v", status)
	}
}

func TestWhisperResult(t *testing.T) {
	response, err := whisper.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	time.Sleep(1 * time.Second)
	result, err := whisper.WhisperResult(response.JobId)
	if err != nil {
		t.Fatalf("Failed to get result: %v", err)
	}
	t.Logf("Result: %v", result)
}

func TestWhisperCancel(t *testing.T) {
	response, err := whisper.WhisperRun(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	cancel, err := whisper.WhisperCancel(response.JobId)
	if err != nil {
		t.Fatalf("Failed to cancel job: %v", err)
	}
	t.Logf("Cancel: %v", cancel)
}

func TestWhisperHealthCheck(t *testing.T) {
	response, err := whisper.WhisperHealthCheck()
	if err != nil {
		t.Fatalf("Failed to health check: %v", err)
	}
	t.Logf("Health check: %v", response)
}

func TestWhisperPurgeQueue(t *testing.T) {
	response, err := whisper.WhisperPurgeQueue()
	if err != nil {
		t.Fatalf("Failed to purge queue: %v", err)
	}
	t.Logf("Purge queue: %v", response)
}
