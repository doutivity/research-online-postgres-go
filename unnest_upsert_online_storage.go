package main

import (
	"context"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"
)

type UnnestUpsertOnlineStorage struct {
	repository *postgresql.Repository
}

func NewUnnestUpsertOnlineStorage(repository *postgresql.Repository) *UnnestUpsertOnlineStorage {
	return &UnnestUpsertOnlineStorage{repository: repository}
}

func (s *UnnestUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	userIDs, timestamps := userOnlinePairsToPgxSlices(pairs)

	return s.repository.Queries().UserOnlineUnnestUpdate(ctx, dbs.UserOnlineUnnestUpdateParams{
		UserIds: userIDs,
		Onlines: timestamps,
	})
}
