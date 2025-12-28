package testcontainers

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
	DSN       string
}

func StartPostgres(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(
		ctx,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"postgres://test:test@%s:%s/testdb?sslmode=disable",
		host,
		port.Port(),
	)

	return &PostgresContainer{
		Container: container,
		DSN:       dsn,
	}, nil
}
