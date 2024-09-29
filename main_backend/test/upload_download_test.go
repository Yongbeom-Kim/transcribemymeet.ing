package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	ctx := context.Background()
	endpoint, teardown := setup(&ctx)
	defer teardown()

	// Start multipart upload
	filePath := "../resources/gettysburg.wav"
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Errorf("Failed to get file info: %v", err)
	}
	fileSizeBytes := fileInfo.Size()
	requestBody := fmt.Sprintf(`{"filename": "%s", "file_size_bytes": %d}`, filepath.Base(filePath), fileSizeBytes)
	resp, err := http.Post(fmt.Sprintf("%s/upload/start-multipart", endpoint), "application/json", strings.NewReader(requestBody))
	if err != nil {
		t.Errorf("Failed to start multipart upload: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to start multipart upload: %v", resp.StatusCode)
	}

	type StartUploadResponse struct {
		UploadID string `json:"upload_id"`
		NumParts int    `json:"num_parts"`
	}

	var response StartUploadResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	partSize := len(fileContent) / response.NumParts
	var parts [][]byte
	for i := 0; i < response.NumParts; i++ {
		start := i * partSize
		end := (i + 1) * partSize
		if i == response.NumParts-1 {
			end = len(fileContent)
		}
		parts = append(parts, fileContent[start:end])
	}

	for part := 0; part < response.NumParts; part++ {
		requestBody := fmt.Sprintf(`{"upload_id": "%s", "part_number": %d}`, response.UploadID, part)
		resp, err := http.Post(fmt.Sprintf("%s/upload/presigned-part-url", endpoint), "application/json", strings.NewReader(requestBody))
		if err != nil {
			t.Fatalf("Failed to get presigned part URL: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
		}

		type PresignedPartURLResponse struct {
			URL string `json:"url"`
		}

		var response PresignedPartURLResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response body: %v", err)
		}

		t.Logf("Uploading part %d to %s", part, response.URL)

		req, err := http.NewRequest("PUT", response.URL, bytes.NewReader(parts[part]))
		if err != nil {
			t.Errorf("Failed to create request for part %d: %v", part, err)
		}
		client := http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			t.Errorf("Failed to upload part %d: %v", part, err)
		}
		defer resp.Body.Close()

		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			t.Errorf("Failed to read part %d response: %v", part, err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to upload part %d: %v", part, resp.StatusCode)
		}
	}

	// Complete multipart upload
	requestBody = fmt.Sprintf(`{"key": "%s", "upload_id": "%s", "num_parts": %d}`, filepath.Base(filePath), response.UploadID, response.NumParts)
	resp, err = http.Post(fmt.Sprintf("%s/upload/complete-multipart", endpoint), "application/json", strings.NewReader(requestBody))
	if err != nil {
		t.Errorf("Failed to complete multipart upload: %v", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Failed to complete multipart upload: %v", resp.StatusCode)
	}

	// Test presigned download URL
	requestBody = fmt.Sprintf(`{"key": "%s"}`, filepath.Base(filePath))
	resp, err = http.Post(fmt.Sprintf("%s/download/presigned-url", endpoint), "application/json", strings.NewReader(requestBody))
	if err != nil {
		t.Errorf("Failed to get presigned download URL: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	type PresignedDownloadURLResponse struct {
		URL string `json:"url"`
	}

	var presignedDownloadURLResponse PresignedDownloadURLResponse
	err = json.Unmarshal(body, &presignedDownloadURLResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	downloadURL := presignedDownloadURLResponse.URL

	t.Logf("Download URL: %s", downloadURL)

	resp, err = http.Get(downloadURL)

	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}
	defer resp.Body.Close()
	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read downloaded content: %v", err)
	}
	// Compare downloaded content with original file
	if !bytes.Equal(downloadedContent, fileContent) {
		t.Fatalf("Downloaded content does not match original file")
	}

}
