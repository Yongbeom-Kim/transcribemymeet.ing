package transcribe

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/whisper"
)

type StartTranscriptionRequest whisper.WhisperInput

type StartTranscriptionResponse struct {
	JobId string `json:"job_id"`
}

func StartTranscription(w http.ResponseWriter, r *http.Request) {
	slog.Info("Starting transcription")
	var reqBody StartTranscriptionRequest
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Unmarshaled request body", "body", reqBody)

	res, err := whisper.WhisperRun(whisper.WhisperInput(reqBody), nil, nil, nil)
	if res == nil {
		http.Error(w, "Received nil response from WhisperRun", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Received response from WhisperRun", "response", res)
	resBody, err := json.Marshal(StartTranscriptionResponse{
		JobId: res.JobId,
	})
	if err != nil {
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
	jobId := r.PathValue("job_id")
	if jobId == "" {
		http.Error(w, "job_id is required", http.StatusBadRequest)
		return
	}

	status, err := whisper.WhisperStatus(jobId)
	if err != nil {
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
	// jobId := r.PathValue("job_id")
	// if jobId == "" {
	// 	http.Error(w, "job_id is required", http.StatusBadRequest)
	// 	return
	// }

	// result, err := whisper.WhisperStatus(jobId)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// if result.Status != runpod.StatusComplete {
	// 	w.WriteHeader(http.StatusAccepted)
	// 	return
	// }

	// if len(result.Output.Segments) == 0 {
	// 	http.Error(w, "Something went wrong. Output is empty", http.StatusInternalServerError)
	// 	return
	// }

	// resBody := GetTranscriptionResultResponse{
	// 	Output: result.Output,
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(resBody)
	panic("Not implemented")
}
