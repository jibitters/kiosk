package containers

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RunPostgres starts a postgres instance as a test container.
func RunPostgres() (testcontainers.Container, int, error) {
	port, _ := nat.NewPort("tcp", "5432")
	request := testcontainers.ContainerRequest{
		Image: "postgres:11",
		Env: map[string]string{
			"POSTGRES_DB":       "kiosk",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort(port),
		AutoRemove:   true,
	}

	container, e := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{ContainerRequest: request, Started: true})
	if e != nil {
		return nil, 0, e
	}

	mappedPort, e := container.MappedPort(context.Background(), port)
	if e != nil {
		return nil, 0, e
	}

	return container, mappedPort.Int(), nil
}

// Stop gets a test container and tries to stop it.
func Stop(container testcontainers.Container) error {
	return container.Terminate(context.Background())
}
