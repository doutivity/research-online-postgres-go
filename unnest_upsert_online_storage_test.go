package main

import (
	"context"
	"testing"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestUnnestUpsertOnlineStorage(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(t, err)
	defer pool.Close()

	storage := NewUnnestUpsertOnlineStorage(postgres.NewDatabase(pool))

	testOnlineStorage(t, storage)
}

func BenchmarkUnnestUpsertOnlineStorage(b *testing.B) {
	if testing.Short() {
		b.Skip()
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(b, err)
	defer pool.Close()

	storage := NewUnnestUpsertOnlineStorage(postgres.NewDatabase(pool))

	benchmarkOnlineStorage(b, storage)
}
