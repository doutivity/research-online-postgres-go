package main

import (
	"context"
	"testing"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"

	"github.com/jackc/pgx/v5"

	"github.com/stretchr/testify/require"
)

func TestUpsertOnlineStorage(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	connection, err := pgx.Connect(ctx, dataSourceName)
	require.NoError(t, err)
	defer connection.Close(ctx)

	storage := NewUpsertOnlineStorage(postgresql.NewSqlcRepository(connection))

	testOnlineStorage(t, storage)
}
