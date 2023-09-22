package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testOnlineStorage(
	t *testing.T,
	storage OnlineStorage,
) {
	t.Helper()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(t, err)
	defer pool.Close()

	var (
		pair1v1 = UserOnlinePair{
			UserID:    1,
			Timestamp: 2e3,
		}
		pair2v1 = UserOnlinePair{
			UserID:    2,
			Timestamp: 3e4,
		}
		pair3v1 = UserOnlinePair{
			UserID:    3,
			Timestamp: 4e5,
		}
		pair4v1 = UserOnlinePair{
			UserID:    4,
			Timestamp: 5e6,
		}
		pair5v1 = UserOnlinePair{
			UserID:    5,
			Timestamp: 6e7,
		}

		pair2v2 = incUserOnlinePair(pair2v1, 10001)
		pair3v2 = incUserOnlinePair(pair3v1, 10002)
	)

	truncateOnline(t, ctx, pool)
	insertOnline(t, ctx, pool, []UserOnlinePair{
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

	expectedOnline(t, ctx, pool, []UserOnlinePair{
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
		batch  = 1000
		online = 1679800725
	)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSourceName)
	require.NoError(b, err)
	defer pool.Close()

	truncateOnline(b, ctx, pool)
	generateOnline(b, ctx, pool, int64(b.N*batch), online)

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
	b.StopTimer()

	expectedOnlineChangedCount(b, ctx, pool, int64(b.N*batch), online)
}

func truncateOnline(t testing.TB, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	// clear
	{
		const (
			// language=PostgreSQL
			query = "TRUNCATE TABLE user_online;"
		)

		_, err := pool.Exec(ctx, query)
		require.NoError(t, err)
	}

}

func insertOnline(t testing.TB, ctx context.Context, pool *pgxpool.Pool, pairs []UserOnlinePair) {
	db := postgres.NewDatabase(pool)

	err := db.WithTransaction(ctx, func(queries *dbs.Queries) error {
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

func generateOnline(t testing.TB, ctx context.Context, pool *pgxpool.Pool, count int64, online int64) {
	t.Helper()

	db := postgres.NewDatabase(pool)
	err := db.Queries().UserOnlineFixtureUpsert(ctx, dbs.UserOnlineFixtureUpsertParams{
		Online: online,
		Count:  count,
	})
	require.NoError(t, err)
}

func expectedOnlineChangedCount(t testing.TB, ctx context.Context, pool *pgxpool.Pool, count int64, online int64) {
	t.Helper()

	db := postgres.NewDatabase(pool)
	row, err := db.Queries().UserOnlineFixtureCount(ctx, online)
	require.NoError(t, err)
	require.Equal(t, count, row.Changed)
}

func expectedOnline(t *testing.T, ctx context.Context, pool *pgxpool.Pool, expectedPairs []UserOnlinePair) {
	t.Helper()

	db := postgres.NewDatabase(pool)
	rows, err := db.Queries().UserOnlineAll(ctx)
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
