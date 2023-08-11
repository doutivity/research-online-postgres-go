package main

import (
	"context"
	"sync/atomic"
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

	truncateOnline(t, ctx, connection)
	insertOnline(t, ctx, connection, []UserOnlinePair{
		pair1v1,
		pair2v1,
		pair3v1,
		pair4v1,
		pair5v1,
	})

	err = storage.BatchStore(ctx, []UserOnlinePair{
		pair2v2,
		pair3v2,
	})
	require.NoError(t, err)

	expectedOnline(t, ctx, connection, []UserOnlinePair{
		pair1v1,
		pair2v2,
		pair3v2,
		pair4v1,
		pair5v1,
	})
}

func benchmarkOnlineStorage(
	b *testing.B,
	storage OnlineStorage,
) {
	b.Helper()

	const (
		batch = 1000
	)

	ctx := context.Background()

	connection, err := pgx.Connect(ctx, dataSourceName)
	require.NoError(b, err)
	defer connection.Close(ctx)

	truncateOnline(b, ctx, connection)
	generateOnline(b, ctx, connection, int64(b.N*batch))

	var (
		startTimestamp = time.Now().Unix()
		counter        = int64(0)
	)

	pairs := make([]UserOnlinePair, batch)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := atomic.AddInt64(&counter, batch)

		for j := 0; j < batch; j++ {
			pairs[j] = UserOnlinePair{
				UserID:    int64(i*batch + j + 1), // 0 .. 99, 100 .. 199
				Timestamp: startTimestamp + index,
			}
		}

		err := storage.BatchStore(ctx, pairs)

		require.NoError(b, err)
	}
}

func truncateOnline(t testing.TB, ctx context.Context, connection *pgx.Conn) {
	t.Helper()

	// clear
	{
		const (
			// language=PostgreSQL
			query = "TRUNCATE TABLE user_online RESTART IDENTITY CASCADE;"
		)

		_, err := connection.Exec(ctx, query)
		require.NoError(t, err)
	}

}

func insertOnline(t testing.TB, ctx context.Context, connection *pgx.Conn, pairs []UserOnlinePair) {
	repository := postgresql.NewSqlcRepository(connection)

	err := repository.WithTransaction(ctx, func(queries *dbs.Queries) error {
		for _, pair := range pairs {
			err := queries.UserOnlineUpsert(ctx, dbs.UserOnlineUpsertParams{
				UserID: pair.UserID,
				Online: pgtype.Timestamp{
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

func generateOnline(t testing.TB, ctx context.Context, connection *pgx.Conn, count int64) {
	const (
		// language=PostgreSQL
		query = `INSERT INTO user_online (user_id, online)
SELECT generate_series,
       to_timestamp(1679800725)
FROM generate_series(1, $1)
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;`
	)

	_, err := connection.Exec(ctx, query, count)
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
