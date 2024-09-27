package gcloud

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is not set", key))
	}
	return value
}

var GetCredentialsFile = sync.OnceValue(func() string {
	return GetEnv("TF_VAR_backend_identity_key")
})

var GetBucket = sync.OnceValue(func() string {
	return GetEnv("TF_VAR_resource_name")
})

func PresignUploadURL(ctx context.Context, key string) (string, error) {

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GetCredentialsFile()))
	if err != nil {
		return "", err
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "PUT",
		Expires: time.Now().Add(15 * time.Minute),
	}

	url, err := client.Bucket(GetBucket()).SignedURL(key, opts)
	if err != nil {
		return "", err
	}

	return url, nil
}

func PresignDownloadURL(ctx context.Context, key string) (string, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GetCredentialsFile()))
	if err != nil {
		return "", err
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	url, err := client.Bucket(GetBucket()).SignedURL(key, opts)
	if err != nil {
		return "", err
	}

	return url, nil
}

func CheckIfObjectExists(ctx context.Context, key string) (bool, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GetCredentialsFile()))
	if err != nil {
		return false, err
	}
	defer client.Close()

	bucket := GetBucket()
	obj := client.Bucket(bucket).Object(key)

	_, err = obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func StartMultipartUpload(ctx context.Context, key string) (uploadID string, err error) {
	randomBytes := make([]byte, 16)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func generatePartKey(key string, uploadID string, partNumber int) (string, error) {
	if partNumber < 0 || partNumber > 31 {
		return "", fmt.Errorf("part number must be between 0 and 31, got %d", partNumber)
	}
	return fmt.Sprintf("%s-%s-part%d", key, uploadID, partNumber), nil
}

func GetUploadPartURL(ctx context.Context, key string, uploadID string, partNumber int) (url string, err error) {
	partKey, err := generatePartKey(key, uploadID, partNumber)
	if err != nil {
		return "", err
	}

	url, err = PresignUploadURL(ctx, partKey)
	if err != nil {
		return "", err
	}

	return url, nil
}

func CompleteMultipartUpload(ctx context.Context, key string, uploadID string, parts int) error {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(GetCredentialsFile()))
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	// 1. Check that all parts are uploaded
	for partNumber := 0; partNumber < parts; partNumber++ {
		partKey, err := generatePartKey(key, uploadID, partNumber)
		if err != nil {
			return fmt.Errorf("failed to generate part key: %v", err)
		}
		exists, err := CheckIfObjectExists(ctx, partKey)
		if err != nil {
			return fmt.Errorf("failed to check if part %d exists: %v", partNumber, err)
		}
		if !exists {
			return fmt.Errorf("part %d not found", partNumber)
		}
	}

	bucket := client.Bucket(GetBucket())

	// 2. Compose all objects into one object
	var sourceObjects []*storage.ObjectHandle
	for partNumber := 0; partNumber < parts; partNumber++ {
		partKey, err := generatePartKey(key, uploadID, partNumber)
		if err != nil {
			return fmt.Errorf("failed to generate part key: %v", err)
		}
		sourceObjects = append(sourceObjects, bucket.Object(partKey))
	}

	compositeObject := bucket.Object(key)
	composer := compositeObject.ComposerFrom(sourceObjects...)
	_, err = composer.Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to compose objects: %v", err)
	}

	// 3. Delete all parts
	// TODO: parallelize this
	for partNumber := 0; partNumber < parts; partNumber++ {
		partKey, err := generatePartKey(key, uploadID, partNumber)
		if err != nil {
			return fmt.Errorf("failed to generate part key: %v", err)
		}
		err = bucket.Object(partKey).Delete(ctx)
		if err != nil {
			fmt.Printf("failed to delete part %s: %v\n", partKey, err)
		}
	}

	return nil
}
