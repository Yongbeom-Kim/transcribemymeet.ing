package whisper

import "github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/runpod"

const (
	StatusQueue    = runpod.StatusQueue
	StatusProgress = runpod.StatusProgress
	StatusComplete = runpod.StatusComplete
	StatusFailed   = runpod.StatusFailed
	StatusCanceled = runpod.StatusCanceled
	StatusTimeout  = runpod.StatusTimeout

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

	// For some reason, this causes the whisper model to not work
	// WhisperTranslationFormatPlainText     = "plain_text"
	// WhisperTranslationFormatFormattedText = "formatted_text"
	// WhisperTranslationFormatSRT           = "srt"
	// WhisperTranslationFormatVTT           = "vtt"
)

type WhisperInput struct {
	AudioURL            string `json:"audio"`
	Model               string `json:"model"`
	TranscriptionFormat string `json:"transcription,omitempty"`
	// For some reason, `translate` causes the whisper model to not work
	// TranslateToEnglish             bool    `json:"translate,omitempty"`
	// TranslationFormat              string  `json:"translation,omitempty"`
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

// Functional Options Pattern

type WhisperInputOption func(*WhisperInput)

func WithModel(model string) WhisperInputOption {
	return func(w *WhisperInput) {
		w.Model = model
	}
}

func WithTranscriptionFormat(format string) WhisperInputOption {
	return func(w *WhisperInput) {
		w.TranscriptionFormat = format
	}
}

// func WithTranslateToEnglish(translate bool) WhisperInputOption {
// 	return func(w *WhisperInput) {
// 		w.TranslateToEnglish = translate
// 	}
// }

// func WithTranslationFormat(format string) WhisperInputOption {
// 	return func(w *WhisperInput) {
// 		w.TranslationFormat = format
// 	}
// }

func WithLanguage(lang string) WhisperInputOption {
	return func(w *WhisperInput) {
		w.Language = lang
	}
}

func WithTemperature(temp float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.Temperature = temp
	}
}

func WithBestOf(best int) WhisperInputOption {
	return func(w *WhisperInput) {
		w.BestOf = best
	}
}

func WithBeamSize(size int) WhisperInputOption {
	return func(w *WhisperInput) {
		w.BeamSize = size
	}
}

func WithPatience(patience float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.Patience = patience
	}
}

func WithLengthPenalty(penalty float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.LengthPenalty = penalty
	}
}

func WithSuppressTokens(tokens string) WhisperInputOption {
	return func(w *WhisperInput) {
		w.SuppressTokens = tokens
	}
}

func WithInitialPrompt(prompt string) WhisperInputOption {
	return func(w *WhisperInput) {
		w.InitialPrompt = prompt
	}
}

func WithConditionOnPreviousText(condition bool) WhisperInputOption {
	return func(w *WhisperInput) {
		w.ConditionOnPreviousText = condition
	}
}

func WithTemperatureIncrementOnFallback(increment float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.TemperatureIncrementOnFallback = increment
	}
}

func WithCompressionRatioThreshold(threshold float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.CompressionRatioThreshold = threshold
	}
}

func WithLogprobThreshold(threshold float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.LogprobThreshold = threshold
	}
}

func WithNoSpeechThreshold(threshold float64) WhisperInputOption {
	return func(w *WhisperInput) {
		w.NoSpeechThreshold = threshold
	}
}

func WithEnableVad(enable bool) WhisperInputOption {
	return func(w *WhisperInput) {
		w.EnableVad = enable
	}
}

func WithWordTimestamps(enable bool) WhisperInputOption {
	return func(w *WhisperInput) {
		w.WordTimestamps = enable
	}
}

func NewWhisperInput(AudioURL string, options ...WhisperInputOption) WhisperInput {
	w := WhisperInput{
		AudioURL:            "",
		Model:               WhisperModelBase,
		TranscriptionFormat: WhisperTranscriptionFormatPlainText,
		// For some reason, `translate` causes the whisper model to not work
		// TranslateToEnglish:  false,
		// TranslationFormat:   WhisperTranscriptionFormatPlainText,
		// Default is None
		// Language:                       nil,
		Temperature: 0,
		BestOf:      5,
		BeamSize:    5,
		// Default is None
		// Patience:                       nil,
		// LengthPenalty:                  nil,
		SuppressTokens:                 "-1",
		InitialPrompt:                  "",
		ConditionOnPreviousText:        false,
		TemperatureIncrementOnFallback: 0.2,
		CompressionRatioThreshold:      2.4,
		LogprobThreshold:               -1,
		NoSpeechThreshold:              0.6,
		EnableVad:                      false,
		WordTimestamps:                 false,
	}
	w.AudioURL = AudioURL
	for _, option := range options {
		option(&w)
	}
	return w
}

type WhisperOutput struct {
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

type WhisperJobStatus struct {
	DelayTime     int    `json:"delayTime,omitempty"`
	ExecutionTime int    `json:"executionTime,omitempty"`
	Status        string `json:"status"`
}
