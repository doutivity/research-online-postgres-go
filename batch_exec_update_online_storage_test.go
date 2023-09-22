package main

import (
	"context"
	"testing"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestBatchExecUpdateOnlineStorage(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(t, err)
	defer pool.Close()

	storage := NewBatchExecUpdateOnlineStorage(postgres.NewDatabase(pool))

	testOnlineStorage(t, storage)
}

func BenchmarkBatchExecUpdateOnlineStorage(b *testing.B) {
	b.Helper()
	if testing.Short() {
		b.Skip()
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(b, err)
	defer pool.Close()

	storage := NewBatchExecUpdateOnlineStorage(postgres.NewDatabase(pool))

	benchmarkOnlineStorage(b, storage)
}
