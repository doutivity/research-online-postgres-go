package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	dataSourceName = "postgresql://yaroslav:AnySecretPassword!!@postgres1:5432/yaaws?sslmode=disable&timezone=UTC"
)

func TestPing(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(t, err)
	defer pool.Close()

	require.NoError(t, pool.Ping(ctx))
}
