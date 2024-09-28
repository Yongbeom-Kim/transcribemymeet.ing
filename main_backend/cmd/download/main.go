package download

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/gcloud"
)

const PRESIGNED_URL_DURATION = 15 * time.Minute

type CreateDownloadURLRequest struct {
	Key string `json:"key"`
}

type CreateDownloadURLResponse struct {
	URL string `json:"url"`
}

func CreateDownloadURL(w http.ResponseWriter, r *http.Request) {
	var req CreateDownloadURLRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := gcloud.PresignDownloadURL(context.Background(), req.Key, PRESIGNED_URL_DURATION)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CreateDownloadURLResponse{URL: url})
}
