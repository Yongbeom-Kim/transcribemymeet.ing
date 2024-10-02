package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestTranscribe(t *testing.T) {
	ctx := context.Background()
	endpoint, teardown := setup(&ctx, t)
	defer teardown()

	AudioFilePath := "https://github.com/runpod-workers/sample-inputs/raw/main/audio/gettysburg.wav"

	// Start transcription
	reqBody := fmt.Sprintf(`{"audio": "%s", "model": "tiny"}`, AudioFilePath)
	resp, err := http.Post(fmt.Sprintf("%s/transcribe/start", endpoint), "application/json", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to start transcription: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code: %v", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if err != nil {
		t.Errorf("Response: %s", string(body))
	}

	t.Logf("Response: %s", string(body))
	type StartResponse struct {
		JobId string `json:"job_id"`
	}

	var startResponse StartResponse
	err = json.Unmarshal(body, &startResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Wait for transcription to complete
	jobId := startResponse.JobId
	t.Logf("JobId: %s", jobId)

	type StatusResponse struct {
		Status string `json:"status"`
	}

Loop:
	for {
		resp, err := http.Get(fmt.Sprintf("%s/transcribe/status/%s", endpoint, jobId))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var statusResponse StatusResponse
		err = json.Unmarshal(body, &statusResponse)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		switch statusResponse.Status {
		case "IN_QUEUE", "IN_PROGRESS":
			t.Logf("Transcription still running")
		case "COMPLETED":
			t.Logf("Transcription completed")
			break Loop
		default:
			t.Fatalf("Invalid transcription status: %s", statusResponse.Status)
		}
		time.Sleep(1 * time.Second)
	}

	// Get transcription results
	resp, err = http.Get(fmt.Sprintf("%s/transcribe/result/%s", endpoint, jobId))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status code: %v", resp.StatusCode)
	}

	t.Logf("Response: %s", string(body))

	var resultResponse struct {
		Output struct {
			Segments []interface{} `json:"segments"`
		} `json:"output"`
	}

	err = json.Unmarshal(body, &resultResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(resultResponse.Output.Segments) == 0 {
		t.Fatalf("No segments found in result: %v", resultResponse)
	}
	t.Logf("Result: %v", resultResponse)
}
