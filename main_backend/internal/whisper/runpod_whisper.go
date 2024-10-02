package whisper

import (
	"encoding/json"
	"fmt"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/runpod"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/utils"
)

var RUNPOD_WHISPER_URL = utils.GetEnvAssert("RUNPOD_WHISPER_URL")

type WhisperSyncRunResponse struct {
	runpod.BaseSyncRunResponse
	Output WhisperOutput `json:"output"`
}

func WhisperRun(input WhisperInput, webhook *runpod.WebHook, policy *runpod.ExecutionPolicy, s3Config *runpod.S3Config) (*runpod.AsyncRunResponse, error) {
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

	response, err := runpod.Run(RUNPOD_WHISPER_URL, runRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func WhisperRunSync(input WhisperInput, webhook *runpod.WebHook, policy *runpod.ExecutionPolicy, s3Config *runpod.S3Config) (*WhisperSyncRunResponse, error) {
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

	response, err := runpod.RunSync(RUNPOD_WHISPER_URL, runRequest)
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

func WhisperStatus(jobId string) (*WhisperJobStatus, error) {
	statusResponse, err := runpod.Status(RUNPOD_WHISPER_URL, jobId)
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

func WhisperResult(jobId string) (*WhisperOutput, error) {
	resultResponse, err := runpod.Status(RUNPOD_WHISPER_URL, jobId)
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

func WhisperCancel(jobId string) (*runpod.CancelResponse, error) {
	cancelResponse, err := runpod.Cancel(RUNPOD_WHISPER_URL, jobId)
	if err != nil {
		return nil, err
	}
	return cancelResponse, nil
}

func WhisperHealthCheck() (*runpod.HealthCheckResponse, error) {
	healthCheckResponse, err := runpod.HealthCheck(RUNPOD_WHISPER_URL)
	if err != nil {
		return nil, err
	}
	return healthCheckResponse, nil
}

func WhisperPurgeQueue() (*runpod.PurgeQueueResponse, error) {
	purgeQueueResponse, err := runpod.PurgeQueue(RUNPOD_WHISPER_URL)
	if err != nil {
		return nil, err
	}
	return purgeQueueResponse, nil
}
