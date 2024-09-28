package main

import (
	"fmt"
	"net/http"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/upload"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/download"
)

func main() {
	http.HandleFunc("POST /upload/start-multipart", upload.StartMultipartUpload)
	http.HandleFunc("POST /upload/presigned-part-url", upload.CreateUploadURL)
	http.HandleFunc("POST /upload/complete-multipart", upload.CompleteMultipartUpload)
	http.HandleFunc("POST /download/presigned-url", download.CreateDownloadURL)

	fmt.Println("Server is starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}

}
