package integration

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setup(ctx *context.Context, t *testing.T) (endpoint string, teardown func()) {
	req := testcontainers.ContainerRequest{
		Image:        "transcribemymeet.ing-backend", // TODO: this should not be hardcoded
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080/tcp"),
		Env: map[string]string{
			"TF_VAR_backend_identity_key": os.Getenv("TF_VAR_backend_identity_key"),
			"TF_VAR_resource_name":        os.Getenv("TF_VAR_resource_name"),
			"RUNPOD_WHISPER_URL":          os.Getenv("RUNPOD_WHISPER_URL"),
			"RUNPOD_API_KEY":              os.Getenv("RUNPOD_API_KEY"),
		},
	}

	container, err := testcontainers.GenericContainer(*ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to start container: %v", err))
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(*ctx, "8080/tcp")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped port: %v", err))
	}

	// Get the IP address
	ip, err := container.Host(*ctx)
	if err != nil {
		panic(fmt.Errorf("failed to get IP address: %v", err))
	}

	return fmt.Sprintf("http://%s:%s", ip, mappedPort.Port()), func() {
		logs, err := container.Logs(*ctx)
		if err != nil {
			t.Fatalf("Failed to get container logs: %v", err)
		}
		defer logs.Close()

		// Read and print the logs
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(logs)
		if err != nil {
			t.Fatalf("Failed to read logs: %v", err)
		}
		logContent := buf.String()

		t.Log("Container Logs:")
		t.Log(logContent)
		container.Terminate(*ctx)
	}
}
