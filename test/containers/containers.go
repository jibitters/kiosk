// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE file.

package containers

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

// ContainersContext is the global context we're going to use to register containers under it.
var ContainersContext = context.Background()

// NewContainer creates a new container based on the given criteria.
func NewContainer(request testcontainers.ContainerRequest) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(
		ContainersContext,
		testcontainers.GenericContainerRequest{ContainerRequest: request, Started: true},
	)
}

// CloseContainer stops the already started container. It's recommended to defer the call to this function.
func CloseContainer(container testcontainers.Container) error {
	return container.Terminate(ContainersContext)
}
