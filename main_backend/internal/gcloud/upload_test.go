package gcloud_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Yongbeom-Kim/transcribemymeet.ing/main_backend/internal/gcloud"
	"google.golang.org/api/option"
)

func generateRandomKey() string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // Handle this error appropriately in your actual code
	}
	return "test-file-" + hex.EncodeToString(randomBytes)
}

func TestPresignUploadURL(t *testing.T) {
	t.Parallel()
	// Set up test environment
	ctx := context.Background()
	key := generateRandomKey()
	testContent := "This is a test file content"
	credentialsFile := gcloud.GetCredentialsFile()

	// 1. Generate presigned URL
	url, err := gcloud.PresignUploadURL(ctx, key)
	if err != nil {
		t.Fatalf("Failed to generate presigned URL: %v", err)
	}

	// 2. Try to upload with presigned URL
	req, err := http.NewRequest("PUT", url, strings.NewReader(testContent))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Upload failed with status code: %d", resp.StatusCode)
	}

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		t.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	defer func() {
		err = storageClient.Bucket(gcloud.GetBucket()).Object(key).Delete(ctx)
		if err != nil {
			t.Logf("Failed to delete test file: %v", err)
		}
	}()

	reader, err := storageClient.Bucket(gcloud.GetBucket()).Object(key).NewReader(ctx)
	if err != nil {
		t.Fatalf("Failed to read uploaded file: %v", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read file content: %v", err)
	}

	if string(content) != testContent {
		t.Fatalf("Uploaded content does not match. Expected: %s, Got: %s", testContent, string(content))
	}
}

func TestPresignUploadDownloadURLFile(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Generate a unique key for the test file
	key := generateRandomKey() + ".wav"

	// Read the file content
	filePath := "../../resources/gettysburg.wav"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Generate presigned upload URL
	uploadURL, err := gcloud.PresignUploadURL(ctx, key)
	if err != nil {
		t.Fatalf("Failed to generate presigned upload URL: %v", err)
	}

	// Upload the file
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
	if err != nil {
		t.Fatalf("Failed to create upload request: %v", err)
	}
	req.Header.Set("Content-Type", "audio/wav")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Upload failed with status code: %d", resp.StatusCode)
	}

	// Generate presigned download URL
	downloadURL, err := gcloud.PresignDownloadURL(ctx, key)
	if err != nil {
		t.Fatalf("Failed to generate presigned download URL: %v", err)
	}

	// Download the file
	downloadResp, err := http.Get(downloadURL)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		t.Fatalf("Download failed with status code: %d", downloadResp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(downloadResp.Body)
	if err != nil {
		t.Fatalf("Failed to read downloaded content: %v", err)
	}

	// Compare the downloaded content with the original file content
	if !bytes.Equal(downloadedContent, fileContent) {
		t.Fatalf("Downloaded content does not match the original file content")
	}

	// Clean up: delete the uploaded file
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(gcloud.GetCredentialsFile()))
	if err != nil {
		t.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	err = storageClient.Bucket(gcloud.GetBucket()).Object(key).Delete(ctx)
	if err != nil {
		t.Logf("Failed to delete test file: %v", err)
	}

}

func TestPresignDownloadURL(t *testing.T) {
	t.Parallel()
	// Set up test environment
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	key := generateRandomKey()
	testContent := fmt.Sprintf("Test content for download %d", time.Now().UnixNano())
	credentialsFile := gcloud.GetCredentialsFile()

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		t.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	// Defer cleanup to ensure it runs even if the test fails
	defer func() {
		err := storageClient.Bucket(gcloud.GetBucket()).Object(key).Delete(ctx)
		if err != nil {
			t.Logf("Failed to delete test file: %v", err)
		}
	}()

	// 1. Manually upload a file
	writer := storageClient.Bucket(gcloud.GetBucket()).Object(key).NewWriter(ctx)
	_, err = writer.Write([]byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Small delay to ensure file is available
	time.Sleep(time.Second)

	// 2. Generate download URL for the file
	downloadURL, err := gcloud.PresignDownloadURL(ctx, key)
	if err != nil {
		t.Fatalf("Failed to generate presigned download URL: %v", err)
	}

	// 3. Verify that the download URL works by downloading the file
	resp, err := http.Get(downloadURL)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Download failed with status code: %d", resp.StatusCode)
	}

	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read downloaded content: %v", err)
	}

	// 4. Compare the file contents
	if string(downloadedContent) != testContent {
		t.Fatalf("Downloaded content does not match. Expected: %s, Got: %s", testContent, string(downloadedContent))
	}
}

func TestCheckIfObjectExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	key := generateRandomKey()
	testContent := "This is a test file for CheckIfObjectExists"

	// Create a storage client
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(gcloud.GetCredentialsFile()))
	if err != nil {
		t.Fatalf("Failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	// Clean up the test object after the test
	defer func() {
		err := storageClient.Bucket(gcloud.GetBucket()).Object(key).Delete(ctx)
		if err != nil {
			t.Logf("Failed to delete test file: %v", err)
		}
	}()

	// Upload a test file
	writer := storageClient.Bucket(gcloud.GetBucket()).Object(key).NewWriter(ctx)
	_, err = writer.Write([]byte(testContent))
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Small delay to ensure file is available
	time.Sleep(time.Second)

	// Check if the object exists
	exists, err := gcloud.CheckIfObjectExists(ctx, key)
	if err != nil {
		t.Fatalf("CheckIfObjectExists failed: %v", err)
	}

	if !exists {
		t.Fatalf("CheckIfObjectExists returned false, expected true")
	}

	// Check for a non-existent object
	nonExistentKey := generateRandomKey()
	exists, err = gcloud.CheckIfObjectExists(ctx, nonExistentKey)
	if err != nil {
		t.Fatalf("CheckIfObjectExists failed for non-existent object: %v", err)
	}

	if exists {
		t.Fatalf("CheckIfObjectExists returned true for non-existent object, expected false")
	}
}

func TestMultipartUploadFile(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		numParts int
	}{
		{"SinglePart", 1},
		{"ThreeParts", 3},
		{"FiveParts", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			key := generateRandomKey()
			filePath := "../../resources/gettysburg.wav"

			// Read the file
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			// Start multipart upload
			uploadID, err := gcloud.StartMultipartUpload(ctx, key)
			if err != nil {
				t.Fatalf("Failed to start multipart upload: %v", err)
			}
			// Split the file into parts
			partSize := len(fileContent) / tc.numParts
			var parts [][]byte
			for i := 0; i < tc.numParts; i++ {
				start := i * partSize
				end := (i + 1) * partSize
				if i == tc.numParts-1 {
					end = len(fileContent)
				}
				parts = append(parts, fileContent[start:end])
			}

			var errChan chan string
			var wg sync.WaitGroup
			for i, part := range parts {
				wg.Add(1)
				go func() {
					defer wg.Done()
					url, err := gcloud.GetUploadPartURL(ctx, key, uploadID, i)
					if err != nil {
						errChan <- fmt.Sprintf("Failed to get upload URL for part %d: %v", i, err)
					}

					req, err := http.NewRequest("PUT", url, bytes.NewReader(part))
					if err != nil {
						errChan <- fmt.Sprintf("Failed to create request for part %d: %v", i, err)
					}

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						errChan <- fmt.Sprintf("Failed to upload part %d: %v", i, err)
					}
					resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						errChan <- fmt.Sprintf("Upload of part %d failed with status code: %d", i, resp.StatusCode)
					}
				}()
			}

			wg.Wait()

			select {
			case err := <-errChan:
				t.Fatalf("Error: %v", err)
			default:
			}

			err = gcloud.CompleteMultipartUpload(ctx, key, uploadID, tc.numParts)
			if err != nil {
				t.Fatalf("Failed to complete multipart upload: %v", err)
			}

			// Deferred cleanup function
			defer func() {
				storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(gcloud.GetCredentialsFile()))
				if err != nil {
					t.Logf("Failed to create storage client for cleanup: %v", err)
					return
				}
				defer storageClient.Close()

				err = storageClient.Bucket(gcloud.GetBucket()).Object(key).Delete(ctx)
				if err != nil {
					t.Logf("Failed to delete test object: %v", err)
				}
			}()

			downloadURL, err := gcloud.PresignDownloadURL(ctx, key)
			if err != nil {
				t.Fatalf("Failed to get download URL: %v", err)
			}
			resp, err := http.Get(downloadURL)
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
		})
	}
}
