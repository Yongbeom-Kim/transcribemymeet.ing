package whisper

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/runpod"
)

type WhisperSyncRunResponse struct {
	runpod.BaseSyncRunResponse
	Output WhisperOutput `json:"output"`
}

type RunpodWhisperClient struct {
	rpclient         *runpod.RunpodClient
	RunpodWhisperURL string
}

var ErrMissingRunpodWhisperURL = fmt.Errorf("runpod whisper URL is required")

func NewRunpodWhisperClient(runpodAPIKey string, runpodWhisperURL string) (*RunpodWhisperClient, error) {
	rpclient, err := runpod.NewRunpodClient(runpodAPIKey)
	if err != nil {
		return nil, err
	}
	if runpodWhisperURL == "" {
		return nil, ErrMissingRunpodWhisperURL
	}

	return &RunpodWhisperClient{
		rpclient:         rpclient,
		RunpodWhisperURL: runpodWhisperURL,
	}, nil
}

func (c *RunpodWhisperClient) Run(input WhisperInput, webhook *runpod.WebHook, policy *runpod.ExecutionPolicy, s3Config *runpod.S3Config) (*runpod.AsyncRunResponse, error) {
	runRequest := runpod.RunRequest{
		Input: input,
	}

	if webhook != nil {
		runRequest.WebHook = *webhook
	}

	if policy != nil {
		runRequest.ExecutionPolicy = *policy
	}

	if s3Config != nil {
		runRequest.S3Config = *s3Config
	}

	slog.Info("Running whisper", "request", runRequest)

	response, err := c.rpclient.Run(c.RunpodWhisperURL, runRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *RunpodWhisperClient) RunSync(input WhisperInput, webhook *runpod.WebHook, policy *runpod.ExecutionPolicy, s3Config *runpod.S3Config) (*WhisperSyncRunResponse, error) {
	runRequest := runpod.RunRequest{
		Input: input,
	}

	if webhook != nil {
		runRequest.WebHook = *webhook
	}

	if policy != nil {
		runRequest.ExecutionPolicy = *policy
	}

	if s3Config != nil {
		runRequest.S3Config = *s3Config
	}

	response, err := c.rpclient.RunSync(c.RunpodWhisperURL, runRequest)
	if err != nil {
		return nil, err
	}

	outputJSON, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response output: %w", err)
	}

	var whisperOutput WhisperSyncRunResponse
	err = json.Unmarshal(outputJSON, &whisperOutput)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transcription: %w", err)
	}

	return &whisperOutput, nil
}

func (c *RunpodWhisperClient) Status(jobId string) (*WhisperJobStatus, error) {
	statusResponse, err := c.rpclient.Status(c.RunpodWhisperURL, jobId)
	if err != nil {
		return nil, err
	}

	outputJSON, err := json.Marshal(statusResponse)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response output: %w", err)
	}

	var whisperStatus WhisperJobStatus
	err = json.Unmarshal(outputJSON, &whisperStatus)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transcription: %w", err)
	}

	return &whisperStatus, nil
}

type ErrJobInProgress struct{}

func (e *ErrJobInProgress) Error() string {
	return "job is currently underway"
}

type ErrJobFailed struct {
	Status string `json:"status"`
}

func (e *ErrJobFailed) Error() string {
	return fmt.Sprintf("job failed with status: %s", e.Status)
}

var JobInProgress = &ErrJobInProgress{}

func (c *RunpodWhisperClient) Result(jobId string) (*WhisperOutput, error) {
	resultResponse, err := c.rpclient.Status(c.RunpodWhisperURL, jobId)
	if err != nil {
		return nil, err
	}

	switch resultResponse.Status {
	case runpod.StatusProgress, runpod.StatusQueue:
		return nil, &ErrJobInProgress{}
	case runpod.StatusComplete:
		break
	case runpod.StatusFailed, runpod.StatusCanceled, runpod.StatusTimeout:
		return nil, &ErrJobFailed{Status: resultResponse.Status}
	default:
		return nil, fmt.Errorf("unknown job status: %s", resultResponse.Status)
	}

	outputJSON, err := json.Marshal(resultResponse.Output)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response output: %w", err)
	}

	var whisperOutput WhisperOutput
	err = json.Unmarshal(outputJSON, &whisperOutput)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transcription: %w", err)
	}
	return &whisperOutput, nil
}

func (c *RunpodWhisperClient) Cancel(jobId string) (*runpod.CancelResponse, error) {
	cancelResponse, err := c.rpclient.Cancel(c.RunpodWhisperURL, jobId)
	if err != nil {
		return nil, err
	}
	return cancelResponse, nil
}

func (c *RunpodWhisperClient) HealthCheck() (*runpod.HealthCheckResponse, error) {
	healthCheckResponse, err := c.rpclient.HealthCheck(c.RunpodWhisperURL)
	if err != nil {
		return nil, err
	}
	return healthCheckResponse, nil
}

func (c *RunpodWhisperClient) PurgeQueue() (*runpod.PurgeQueueResponse, error) {
	purgeQueueResponse, err := c.rpclient.PurgeQueue(c.RunpodWhisperURL)
	if err != nil {
		return nil, err
	}
	return purgeQueueResponse, nil
}
