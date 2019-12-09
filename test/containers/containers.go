// Copyright 2019 The JIBit Team. All rights reserved.
// Use of this source code is governed by an Apache Style license that can be found in the LICENSE.md file.

package containers

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
)

// ContainersContext is the global context we're going to use to register containers under it.
var ContainersContext = context.Background()

// NewContainer creates a new container based on the given criteria. If we manage to successfully create and
// start a new container, then the first returned value would be the started container. Otherwise, you should
// checkout the error instance for more information on what exactly went wrong.
func NewContainer(request testcontainers.ContainerRequest) (testcontainers.Container, error) {
	return testcontainers.GenericContainer(
		ContainersContext,
		testcontainers.GenericContainerRequest{ContainerRequest: request, Started: true},
	)
}

// CloseContainer stops the already started container. If we fail to stop the given container, then we would return
// an error. Also, it's recommended to defer the call to this function.
func CloseContainer(container testcontainers.Container) error {
	return container.Terminate(ContainersContext)
}
