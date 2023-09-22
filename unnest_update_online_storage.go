package main

import (
	"context"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"
)

type BatchUpdateOnlineStorage struct {
	db *postgres.Database
}

func NewUnnestUpdateOnlineStorage(db *postgres.Database) *BatchUpdateOnlineStorage {
	return &BatchUpdateOnlineStorage{db: db}
}

func (s *BatchUpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	userIDs, timestamps := userOnlinePairsToPgxSlices(pairs)

	return s.db.Queries().UserOnlineUnnestUpdate(ctx, dbs.UserOnlineUnnestUpdateParams{
		UserIds: userIDs,
		Onlines: timestamps,
	})
}
