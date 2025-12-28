package repository

import (
	"context"
	"os"
	"testing"

	tc "kenjix.com/test/testcontainers"
)

var testDSN string
var ctx = context.Background()

func TestMain(m *testing.M) {
	pg, err := tc.StartPostgres(ctx)
	if err != nil {
		panic(err)
	}
	defer pg.Container.Terminate(ctx)

	jdbcURL := "jdbc:" + pg.DSN

	if err := tc.RunFlyway(ctx, jdbcURL); err != nil {
		panic(err)
	}

	testDSN = pg.DSN

	code := m.Run()
	os.Exit(code)
}
