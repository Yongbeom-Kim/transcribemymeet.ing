package whisper_test

import (
	"os"
	"testing"
	"time"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/whisper"
)

var input = whisper.NewWhisperInput(
	"https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav",
	whisper.WithModel(whisper.WhisperModelTiny),
)

func getRunpodWhisperClient() (*whisper.RunpodWhisperClient, error) {
	c, err := whisper.NewRunpodWhisperClient(os.Getenv("RUNPOD_API_KEY"), os.Getenv("RUNPOD_WHISPER_URL"))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func TestWhisperRun(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.Run(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperRunSync(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.RunSync(input, nil, nil, nil)
	if err != nil {
		t.Errorf("Failed to run job: %v", err)
	}
	t.Logf("Response: %v", response)
}

func TestWhisperStatus(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.Run(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	status, err := c.Status(response.JobId)
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
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.Run(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	time.Sleep(2 * time.Second)
	result, err := c.Result(response.JobId)
	if err != nil {
		t.Fatalf("Failed to get result: %v", err)
	}
	t.Logf("Result: %v", result)
}

func TestWhisperCancel(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.Run(input, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to run job: %v", err)
	}
	cancel, err := c.Cancel(response.JobId)
	if err != nil {
		t.Fatalf("Failed to cancel job: %v", err)
	}
	t.Logf("Cancel: %v", cancel)
}

func TestWhisperHealthCheck(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.HealthCheck()
	if err != nil {
		t.Fatalf("Failed to health check: %v", err)
	}
	t.Logf("Health check: %v", response)
}

func TestWhisperPurgeQueue(t *testing.T) {
	c, err := getRunpodWhisperClient()
	if err != nil {
		t.Fatalf("Failed to get runpod whisper client: %v", err)
		return
	}

	response, err := c.PurgeQueue()
	if err != nil {
		t.Fatalf("Failed to purge queue: %v", err)
	}
	t.Logf("Purge queue: %v", response)
}
