// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package nats

import (
	"strconv"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jibitters/kiosk/internal/app/kiosk/configuration"
	"github.com/jibitters/kiosk/test/containers"
	"github.com/testcontainers/testcontainers-go"
)

func runNatsContainer() (testcontainers.Container, int, error) {
	containerPort, err := nat.NewPort("tcp", "4222")
	if err != nil {
		return nil, 0, err
	}

	request := testcontainers.ContainerRequest{
		Image:        "nats:2",
		ExposedPorts: []string{"4222/tcp"},
	}

	container, err := containers.NewContainer(request)
	if err != nil {
		return nil, 0, err
	}

	mappedPort, err := container.MappedPort(containers.ContainersContext, containerPort)
	if err != nil {
		return nil, 0, err
	}

	return container, mappedPort.Int(), nil
}

func TestConnectToNats(t *testing.T) {
	container, port, err := runNatsContainer()
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer containers.CloseContainer(container)

	config := &configuration.Config{Nats: configuration.NatsConfig{
		Addresses: []string{"nats://localhost:" + strconv.Itoa(port)},
	}}

	time.Sleep(time.Second)
	nats, err := ConnectToNats(config)
	if err != nil {
		t.Logf("Error : %v", err)
		t.FailNow()
	}
	defer nats.Close()
}

func TestConnectToNats_Error(t *testing.T) {
	config := &configuration.Config{Nats: configuration.NatsConfig{
		Addresses: []string{"nats://localhost:" + strconv.Itoa(4222+4222)},
	}}

	_, err := ConnectToNats(config)
	if err == nil {
		t.Logf("Expected error here : %v", err)
		t.FailNow()
	}
}
