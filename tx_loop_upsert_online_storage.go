package main

import (
	"context"
	"time"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type TxLoopUpsertOnlineStorage struct {
	db *postgres.Database
}

func NewTxLoopUpsertOnlineStorage(db *postgres.Database) *TxLoopUpsertOnlineStorage {
	return &TxLoopUpsertOnlineStorage{db: db}
}

func (s *TxLoopUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	return s.db.WithTransaction(ctx, func(queries *dbs.Queries) error {
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
}
