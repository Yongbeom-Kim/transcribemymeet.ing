package runpod

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

const (
	StatusQueue    = "IN_QUEUE"    // Job is waiting in the endpoint queue for an available worker
	StatusProgress = "IN_PROGRESS" // Job is actively being processed by a worker
	StatusComplete = "COMPLETED"   // Job has successfully finished processing and returned a result
	StatusFailed   = "FAILED"      // Job encountered an error during execution
	StatusCanceled = "CANCELLED"   // Job was manually cancelled before completion
	StatusTimeout  = "TIMED_OUT"   // Job expired before processing or worker failed to report result in time
)

var RUNPOD_API_KEY = sync.OnceValue(func() string {
	key := os.Getenv("RUNPOD_API_KEY")
	if key == "" {
		panic("RUNPOD_API_KEY environment variable is not set")
	}
	return key
})

type WebHook string
type ExecutionPolicy struct {
	Timeout    int `json:"executionTimeout"`
	Priority   int `json:"priority"`
	TimeToLive int `json:"ttl"`
}

type S3Config struct {
	AccessId     string `json:"accessId"`
	AccessSecret string `json:"accessSecret"`
	BucketName   string `json:"bucketName"`
	EndpointURL  string `json:"endpointUrl"`
}

type RunRequest struct {
	Input           interface{}     `json:"input"`
	WebHook         WebHook         `json:"webhook,omitempty"`
	ExecutionPolicy ExecutionPolicy `json:"policy,omitempty"`
	S3Config        S3Config        `json:"s3Config,omitempty"`
}

type AsyncRunResponse struct {
	JobId  string `json:"id"`
	Status string `json:"status"`
}

type BaseSyncRunResponse struct {
	DelayTime     int    `json:"delayTime"`
	ExecutionTime int    `json:"executionTime"`
	JobId         string `json:"id"`
	Status        string `json:"status"`
}

type SyncRunResponse struct {
	BaseSyncRunResponse
	Output map[string]interface{} `json:"output"`
}

func Run(workerURL string, runRequest RunRequest) (*AsyncRunResponse, error) {
	slog.Info("Running job", "workerURL", workerURL)
	jsonData, err := json.Marshal(runRequest)
	if err != nil {
		slog.Error("Error marshalling run request", "error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/run", workerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating run request", "error", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error running job", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading run response", "error", err)
		return nil, err
	}

	var runResponse AsyncRunResponse
	err = json.Unmarshal(body, &runResponse)
	if err != nil {
		slog.Error("Error unmarshalling run response", "error", err)
		return nil, err
	}

	slog.Info("Run response", "response", runResponse)

	return &runResponse, nil
}

func RunSync(workerURL string, runRequest RunRequest) (*SyncRunResponse, error) {
	slog.Info("Running job synchronously", "workerURL", workerURL)
	jsonData, err := json.Marshal(runRequest)
	if err != nil {
		slog.Error("Error marshalling run request", "error", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/runsync", workerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating run request", "error", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error running job", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading run response", "error", err)
		return nil, err
	}

	var runResponse SyncRunResponse
	err = json.Unmarshal(body, &runResponse)
	if err != nil {
		slog.Error("Error unmarshalling run response", "error", err)
		return nil, err
	}

	slog.Info("Run response", "response", runResponse)

	return &runResponse, nil
}

// func Stream(workerURL string, jobId string) (string, error) {

// }

type StatusResponse struct {
	DelayTime     int         `json:"delayTime,omitempty"`
	ExecutionTime int         `json:"executionTime,omitempty"`
	JobId         string      `json:"id"`
	Status        string      `json:"status"`
	Output        interface{} `json:"output,omitempty"`
}

var ErrEmtpyStatus = errors.New("empty status received")

func Status(workerURL string, jobId string) (*StatusResponse, error) {
	slog.Info("Getting status", "workerURL", workerURL, "jobId", jobId)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status/%s", workerURL, jobId), nil)
	if err != nil {
		slog.Error("Error creating status request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error getting status", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading status response", "error", err)
		return nil, err
	}

	var statusResponse StatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		slog.Error("Error unmarshalling status response", "error", err)
		return nil, err
	}
	// Check if statusResponse.Status is one of the defined constants
	switch statusResponse.Status {
	case StatusQueue, StatusProgress, StatusComplete, StatusFailed, StatusCanceled, StatusTimeout:
		// Status is valid
	case "":
		return nil, ErrEmtpyStatus
	default:
		// Status is not one of the expected values
		slog.Error("Unexpected status received", "status", statusResponse.Status)
		return nil, fmt.Errorf("unexpected status received: %s", statusResponse.Status)
	}

	slog.Info("Status response", "response", statusResponse)

	return &statusResponse, nil
}

type CancelResponse struct {
	JobId  string `json:"id"`
	Status string `json:"status"`
}

func Cancel(workerURL string, jobId string) (*CancelResponse, error) {
	slog.Info("Cancelling job", "workerURL", workerURL, "jobId", jobId)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cancel/%s", workerURL, jobId), nil)
	if err != nil {
		slog.Error("Error creating cancel request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error cancelling job", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading cancel response", "error", err)
		return nil, err
	}

	var cancelResponse CancelResponse
	err = json.Unmarshal(body, &cancelResponse)
	if err != nil {
		slog.Error("Error unmarshalling cancel response", "error", err)
		return nil, err
	}

	slog.Info("Cancel response", "response", cancelResponse)

	return &cancelResponse, nil
}

type HealthCheckResponse struct {
	Jobs struct {
		Completed  int `json:"completed"`
		Failed     int `json:"failed"`
		InProgress int `json:"inProgress"`
		InQueue    int `json:"inQueue"`
		Retried    int `json:"retried"`
	} `json:"jobs"`
	Workers struct {
		Idle         int `json:"idle"`
		Initializing int `json:"initializing"`
		Ready        int `json:"ready"`
		Running      int `json:"running"`
		Throttled    int `json:"throttled"`
		Unhealthy    int `json:"unhealthy"`
	} `json:"workers"`
}

func HealthCheck(workerURL string) (*HealthCheckResponse, error) {
	slog.Info("Checking health", "workerURL", workerURL)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/health", workerURL), nil)
	if err != nil {
		slog.Error("Error creating health check request", "error", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error checking health", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading health response", "error", err)
		return nil, err
	}

	var healthCheckResponse HealthCheckResponse
	err = json.Unmarshal(body, &healthCheckResponse)
	if err != nil {
		slog.Error("Error unmarshalling health response", "error", err)
		return nil, err
	}

	slog.Info("Health check response", "response", healthCheckResponse)

	return &healthCheckResponse, nil
}

type PurgeQueueResponse struct {
	JobsRemoved int    `json:"removed"`
	Status      string `json:"status"`
}

func PurgeQueue(workerURL string) (*PurgeQueueResponse, error) {
	slog.Info("Purging queue", "workerURL", workerURL)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/purge-queue", workerURL), nil)
	if err != nil {
		slog.Error("Error creating purge queue request", "error", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", RUNPOD_API_KEY()))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Error purging queue", "error", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading purge queue response", "error", err)
		return nil, err
	}

	var purgeQueueResponse PurgeQueueResponse
	err = json.Unmarshal(body, &purgeQueueResponse)
	if err != nil {
		slog.Error("Error unmarshalling purge queue response", "error", err)
		return nil, err
	}

	slog.Info("Purge queue response", "response", purgeQueueResponse)

	return &purgeQueueResponse, nil
}
