package runpod

import (
	"encoding/json"
	"fmt"
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

const (
	WhisperModelTiny    = "tiny"
	WhisperModelBase    = "base"
	WhisperModelSmall   = "small"
	WhisperModelMedium  = "medium"
	WhisperModelLargeV1 = "large-v1"
	WhisperModelLargeV2 = "large-v2"
	WhisperModelLargeV3 = "large-v3"

	WhisperTranscriptionFormatPlainText     = "plain_text"
	WhisperTranscriptionFormatFormattedText = "formatted_text"
	WhisperTranscriptionFormatSRT           = "srt"
	WhisperTranscriptionFormatVTT           = "vtt"

	WhisperTranslationFormatPlainText     = "plain_text"
	WhisperTranslationFormatFormattedText = "formatted_text"
	WhisperTranslationFormatSRT           = "srt"
	WhisperTranslationFormatVTT           = "vtt"
)

type WhisperModelInput struct {
	Audio                          string  `json:"audio,omitempty"`
	AudioBase64                    string  `json:"audio_base64,omitempty"`
	Model                          string  `json:"model,omitempty"`
	Transcription                  string  `json:"transcription,omitempty"`
	Translate                      bool    `json:"translate,omitempty"`
	Translation                    string  `json:"translation,omitempty"`
	Language                       string  `json:"language,omitempty"`
	Temperature                    float64 `json:"temperature,omitempty"`
	BestOf                         int     `json:"best_of,omitempty"`
	BeamSize                       int     `json:"beam_size,omitempty"`
	Patience                       float64 `json:"patience,omitempty"`
	LengthPenalty                  float64 `json:"length_penalty,omitempty"`
	SuppressTokens                 string  `json:"suppress_tokens,omitempty"`
	InitialPrompt                  string  `json:"initial_prompt,omitempty"`
	ConditionOnPreviousText        bool    `json:"condition_on_previous_text,omitempty"`
	TemperatureIncrementOnFallback float64 `json:"temperature_increment_on_fallback,omitempty"`
	CompressionRatioThreshold      float64 `json:"compression_ratio_threshold,omitempty"`
	LogprobThreshold               float64 `json:"logprob_threshold,omitempty"`
	NoSpeechThreshold              float64 `json:"no_speech_threshold,omitempty"`
	EnableVad                      bool    `json:"enable_vad,omitempty"`
	WordTimestamps                 bool    `json:"word_timestamps,omitempty"`
}

type WhisperModelOutput struct {
	Segments []struct {
		ID               int     `json:"id"`
		Seek             int     `json:"seek"`
		Start            float64 `json:"start"`
		End              float64 `json:"end"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens"`
		Temperature      float64 `json:"temperature"`
		AvgLogprob       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	} `json:"segments"`
	DetectedLanguage string      `json:"detected_language"`
	Transcription    string      `json:"transcription"`
	Translation      interface{} `json:"translation"`
	Device           string      `json:"device"`
	Model            string      `json:"model"`
	TranslationTime  float64     `json:"translation_time"`
}

type WhisperSyncRunResponse struct {
	BaseSyncRunResponse
	Output WhisperModelOutput `json:"output"`
}

type WhisperStatusResponse struct {
	BaseStatusResponse
	Output WhisperModelOutput `json:"output,omitempty"`
}

func WhisperRun(input WhisperModelInput, webhook *WebHook, policy *ExecutionPolicy, s3Config *S3Config) (*AsyncRunResponse, error) {
	runRequest := RunRequest{
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

	response, err := Run(RUNPOD_WHISPER_URL(), runRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func WhisperRunSync(input WhisperModelInput, webhook *WebHook, policy *ExecutionPolicy, s3Config *S3Config) (*WhisperSyncRunResponse, error) {
	runRequest := RunRequest{
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

	response, err := RunSync(RUNPOD_WHISPER_URL(), runRequest)
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

func WhisperStatus(jobId string) (*WhisperStatusResponse, error) {
	statusResponse, err := Status(RUNPOD_WHISPER_URL(), jobId)
	if err != nil {
		return nil, err
	}

	outputJSON, err := json.Marshal(statusResponse)
	if err != nil {
		return nil, fmt.Errorf("error marshaling response output: %w", err)
	}

	var whisperStatus WhisperStatusResponse
	err = json.Unmarshal(outputJSON, &whisperStatus)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling transcription: %w", err)
	}

	return &whisperStatus, nil
}

func WhisperCancel(jobId string) (*CancelResponse, error) {
	cancelResponse, err := Cancel(RUNPOD_WHISPER_URL(), jobId)
	if err != nil {
		return nil, err
	}
	return cancelResponse, nil
}

func WhisperHealthCheck() (*HealthCheckResponse, error) {
	healthCheckResponse, err := HealthCheck(RUNPOD_WHISPER_URL())
	if err != nil {
		return nil, err
	}
	return healthCheckResponse, nil
}

func WhisperPurgeQueue() (*PurgeQueueResponse, error) {
	purgeQueueResponse, err := PurgeQueue(RUNPOD_WHISPER_URL())
	if err != nil {
		return nil, err
	}
	return purgeQueueResponse, nil
}
