package main

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
)

func StartContainerAndGetHost(ctx context.Context, req testcontainers.ContainerRequest, mappedPort nat.Port) (testcontainers.Container, string) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	ip, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err := container.MappedPort(ctx, mappedPort)
	if err != nil {
		panic(err)
	}
	return container, fmt.Sprintf("%s:%s", ip, port.Port())
}
