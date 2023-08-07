package main

import (
	"context"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"
)

type BatchUpsertOnlineStorage struct {
	repository *postgresql.Repository
}

func NewBatchUpsertOnlineStorage(repository *postgresql.Repository) *BatchUpsertOnlineStorage {
	return &BatchUpsertOnlineStorage{repository: repository}
}

func (s *BatchUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	userIDs, timestamps := userOnlinePairsToPgxSlices(pairs)

	return s.repository.Queries().UserOnlineBatchUpsert(ctx, dbs.UserOnlineBatchUpsertParams{
		UserIds: userIDs,
		Onlines: timestamps,
	})
}
