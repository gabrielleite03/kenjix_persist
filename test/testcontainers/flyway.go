package testcontainers

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunFlyway(ctx context.Context, jdbcURL string) error {
	req := testcontainers.ContainerRequest{
		Image: "flyway/flyway:9",
		Cmd:   []string{"migrate"},
		Env: map[string]string{
			"FLYWAY_URL":      jdbcURL,
			"FLYWAY_USER":     "test",
			"FLYWAY_PASSWORD": "test",
		},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount("./db/migration", "/flyway/sql"),
		),
		WaitingFor: wait.ForExit(),
	}

	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return err
	}

	return container.Terminate(ctx)
}
