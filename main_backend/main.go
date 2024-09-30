package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/download"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/cmd/upload"
)

func main() {
	http.HandleFunc("POST /upload/start-multipart", upload.StartMultipartUpload)
	http.HandleFunc("POST /upload/presigned-part-url", upload.CreateUploadURL)
	http.HandleFunc("POST /upload/complete-multipart", upload.CompleteMultipartUpload)
	http.HandleFunc("POST /download/presigned-url", download.CreateDownloadURL)

	port := os.Getenv("PORT")
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
