package main

import (
	"context"
	"testing"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/stretchr/testify/require"
)

func testOnlineStorage(
	t *testing.T,
	storage OnlineStorage,
) {
	t.Helper()
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	connection, err := pgx.Connect(ctx, dataSourceName)
	require.NoError(t, err)
	defer connection.Close(ctx)

	var (
		pair1v1 = UserOnlinePair{
			UserID:    1,
			Timestamp: 1679800725,
		}
		pair2v1 = UserOnlinePair{
			UserID:    2,
			Timestamp: 1679800730,
		}
		pair3v1 = UserOnlinePair{
			UserID:    3,
			Timestamp: 1679800735,
		}
		pair4v1 = UserOnlinePair{
			UserID:    4,
			Timestamp: 1679800740,
		}
		pair5v1 = UserOnlinePair{
			UserID:    5,
			Timestamp: 1679800745,
		}

		pair2v2 = incUserOnlinePair(pair2v1, 10001)
		pair3v2 = incUserOnlinePair(pair3v1, 10002)
	)

	resetOnline(t, ctx, connection, []UserOnlinePair{
		pair1v1,
		pair2v1,
		pair3v1,
		pair4v1,
		pair5v1,
	})

	storage.BatchStore(ctx, []UserOnlinePair{
		pair2v2,
		pair3v2,
	})

	expectedOnline(t, ctx, connection, []UserOnlinePair{
		pair1v1,
		pair2v2,
		pair3v2,
		pair4v1,
		pair5v1,
	})
}

func resetOnline(t *testing.T, ctx context.Context, connection *pgx.Conn, pairs []UserOnlinePair) {
	t.Helper()

	// clear
	{
		const (
			// language =PostgreSQL
			query = "TRUNCATE TABLE user_online RESTART IDENTITY CASCADE"
		)

		_, err := connection.Exec(ctx, query)
		require.NoError(t, err)
	}

	repository := postgresql.NewSqlcRepository(connection)

	err := repository.WithTransaction(ctx, func(queries *dbs.Queries) error {
		for _, pair := range pairs {
			err := queries.UserOnlineNew(ctx, dbs.UserOnlineNewParams{
				UserID: pair.UserID,
				Online: pgtype.Timestamptz{
					Time:  time.Unix(pair.Timestamp, 0).UTC(),
					Valid: true,
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	require.NoError(t, err)
}

func expectedOnline(t *testing.T, ctx context.Context, connection *pgx.Conn, expectedPairs []UserOnlinePair) {
	t.Helper()

	repository := postgresql.NewSqlcRepository(connection)
	rows, err := repository.Queries().UserOnlineAll(ctx)
	require.NoError(t, err)

	actualPairs := make([]UserOnlinePair, len(rows))
	for i, row := range rows {
		actualPairs[i] = UserOnlinePair{
			UserID:    row.UserID,
			Timestamp: row.Online.Time.Unix(),
		}
	}
	require.Equal(t, expectedPairs, actualPairs)
}
