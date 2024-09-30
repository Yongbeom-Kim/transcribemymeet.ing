package upload

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/gcloud"
	"github.com/google/uuid"
)

const PRESIGNED_URL_DURATION = 15 * time.Minute

type StartUploadRequest struct {
	Filename      string `json:"filename"`
	FileSizeBytes int    `json:"file_size_bytes"`
}

type StartUploadResponse struct {
	UploadID string `json:"upload_id"`
	NumParts int    `json:"num_parts"`
}

func StartMultipartUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StartUploadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadID := uuid.New().String()
	numParts := calculateNumParts(req.FileSizeBytes)

	resp := StartUploadResponse{
		UploadID: uploadID,
		NumParts: numParts,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func calculateNumParts(fileSizeBytes int) int {
	partSize := 64 * 1024 * 1024 // 64 MB
	numParts := fileSizeBytes / partSize
	if numParts == 0 {
		return 1
	}
	if numParts > 32 {
		return 32
	}
	return numParts
}

type CreateUploadURLRequest struct {
	UploadID   string `json:"upload_id"`
	PartNumber int    `json:"part_number"`
}

type CreateUploadURLResponse struct {
	URL string `json:"url"`
}

func CreateUploadURL(w http.ResponseWriter, r *http.Request) {
	var req CreateUploadURLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UploadID == "" || req.PartNumber < 0 || req.PartNumber >= 32 { // TODO: make this a constant
		http.Error(w, "Invalid upload_id or part_number", http.StatusBadRequest)
		return
	}

	url, err := gcloud.GetUploadPartURL(context.Background(), req.UploadID, req.PartNumber, PRESIGNED_URL_DURATION)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreateUploadURLResponse{URL: url})
}

type CompleteMultipartUploadRequest struct {
	Key      string `json:"key"`
	UploadID string `json:"upload_id"`
	NumParts int    `json:"num_parts"`
}

func CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	var req CompleteMultipartUploadRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = gcloud.CompleteMultipartUpload(context.Background(), req.Key, req.UploadID, req.NumParts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
