package integration

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setup(ctx *context.Context) (endpoint string, teardown func()) {
	req := testcontainers.ContainerRequest{
		Image:        "transcribemymeet.ing-backend", // TODO: this should not be hardcoded
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForListeningPort("8080/tcp"),
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
		container.Terminate(*ctx)
	}
}
