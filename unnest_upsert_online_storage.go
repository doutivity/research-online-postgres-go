package main

import (
	"context"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"
)

type UnnestUpsertOnlineStorage struct {
	db *postgres.Database
}

func NewUnnestUpsertOnlineStorage(db *postgres.Database) *UnnestUpsertOnlineStorage {
	return &UnnestUpsertOnlineStorage{db: db}
}

func (s *UnnestUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	userIDs, timestamps := userOnlinePairsToPgxSlices(pairs)

	return s.db.Queries().UserOnlineUnnestUpsert(ctx, dbs.UserOnlineUnnestUpsertParams{
		UserIds: userIDs,
		Onlines: timestamps,
	})
}
