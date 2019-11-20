package containers

import (
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewContainer(t *testing.T) {
	containerPort, err := nat.NewPort("tcp", "5432")
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	request := testcontainers.ContainerRequest{
		Image:        "postgres:11",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort(containerPort),
	}

	container, err := NewContainer(request)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer CloseContainer(container)

	ports, err := container.Ports(ContainersContext)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}

	bindings := ports[containerPort]
	if len(bindings) != 1 {
		t.Logf("Actual value: %v Expected value: 1", len(bindings))
		t.FailNow()
	}
}
