package main

import (
	"context"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"
)

type BatchUpdateOnlineStorage struct {
	repository *postgresql.Repository
}

func NewUnnestUpdateOnlineStorage(repository *postgresql.Repository) *BatchUpdateOnlineStorage {
	return &BatchUpdateOnlineStorage{repository: repository}
}

func (s *BatchUpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	userIDs, timestamps := userOnlinePairsToPgxSlices(pairs)

	return s.repository.Queries().UserOnlineUnnestUpdate(ctx, dbs.UserOnlineUnnestUpdateParams{
		UserIds: userIDs,
		Onlines: timestamps,
	})
}
