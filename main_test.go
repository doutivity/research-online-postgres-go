package main

import (
	"context"
	"testing"
	"time"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	dataSourceName = "postgresql://yaroslav:AnySecretPassword!!@postgres1:5432/yaaws?sslmode=disable&timezone=UTC"
)

func TestPing(t *testing.T) {
	connection, err := pgx.Connect(context.Background(), dataSourceName)
	require.NoError(t, err)
	defer connection.Close(context.Background())

	require.NoError(t, connection.Ping(context.Background()))

	queries := dbs.New(connection)
	err = queries.UserOnlineNew(context.Background(), dbs.UserOnlineNewParams{
		UserOnline: 1,
		Online: pgtype.Timestamptz{
			Time:             time.Now(),
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	require.NoError(t, err)
}
