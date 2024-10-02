package transcribe

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/whisper"
)

type StartTranscriptionRequest whisper.WhisperInput

type StartTranscriptionResponse struct {
	JobId string `json:"job_id"`
}

func getRunpodWhisperClient() (*whisper.RunpodWhisperClient, error) {
	return whisper.NewRunpodWhisperClient(os.Getenv("RUNPOD_API_KEY"), os.Getenv("RUNPOD_WHISPER_URL"))
}

func StartTranscription(w http.ResponseWriter, r *http.Request) {
	slog.Info("Starting transcription")
	whisperClient, err := getRunpodWhisperClient()
	if err != nil {
		slog.Error("Failed to get RunpodWhisperClient", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var reqBody StartTranscriptionRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		slog.Error("Failed to unmarshal request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Unmarshaled request body", "body", reqBody)

	res, err := whisperClient.Run(whisper.WhisperInput(reqBody), nil, nil, nil)
	if res == nil {
		slog.Error("Received nil response from WhisperRun")
		http.Error(w, "Received nil response from WhisperRun", http.StatusInternalServerError)
		return
	}
	if err != nil {
		slog.Error("Failed to run Whisper", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Received response from WhisperRun", "response", res)
	resBody, err := json.Marshal(StartTranscriptionResponse{
		JobId: res.JobId,
	})
	if err != nil {
		slog.Error("Failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Writing response to client", "response", resBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

type GetTranscriptionStatusResponse struct {
	Status        string `json:"status"`
	DelayTime     int    `json:"delay_time,omitempty"`
	ExecutionTime int    `json:"execution_time,omitempty"`
}

func GetTranscriptionStatus(w http.ResponseWriter, r *http.Request) {
	whisperClient, err := getRunpodWhisperClient()
	if err != nil {
		slog.Error("Failed to get RunpodWhisperClient", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobId := r.PathValue("job_id")
	if jobId == "" {
		slog.Error("job_id is required")
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}

	status, err := whisperClient.Status(jobId)
	if err != nil {
		slog.Error("Failed to get status", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := GetTranscriptionStatusResponse{
		Status:        status.Status,
		DelayTime:     status.DelayTime,
		ExecutionTime: status.ExecutionTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}

type GetTranscriptionResultResponse struct {
	Output whisper.WhisperOutput `json:"output"`
}

func GetTranscriptionResult(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting transcription result")
	whisperClient, err := getRunpodWhisperClient()
	if err != nil {
		slog.Error("Failed to get RunpodWhisperClient", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jobId := r.PathValue("job_id")
	if jobId == "" {
		slog.Error("job_id is required")
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}

	result, err := whisperClient.Result(jobId)
	if err != nil {
		slog.Error("Failed to get result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody := GetTranscriptionResultResponse{
		Output: *result,
	}

	slog.Info("Writing transcription result to client", "response", resBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resBody)
}
