package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/download"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/transcribe"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/upload"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/utils"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	http.HandleFunc("POST /upload/start-multipart", upload.StartMultipartUpload)
	http.HandleFunc("POST /upload/presigned-part-url", upload.CreateUploadURL)
	http.HandleFunc("POST /upload/complete-multipart", upload.CompleteMultipartUpload)
	http.HandleFunc("POST /download/presigned-url", download.CreateDownloadURL)
	http.HandleFunc("POST /transcribe/start", transcribe.StartTranscription)
	http.HandleFunc("GET /transcribe/status/{job_id}", transcribe.GetTranscriptionStatus)
	http.HandleFunc("GET /transcribe/result/{job_id}", transcribe.GetTranscriptionResult)

	port := utils.GetEnvAssert("PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Server is starting on port %d...\n", portInt)
	err = http.ListenAndServe(fmt.Sprintf(":%d", portInt), nil)
	if err != nil {
		fmt.Println(err)
	}

}
