package main

import (
	"context"
	"time"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type TxLoopUpdateOnlineStorage struct {
	db *postgres.Database
}

func NewTxLoopUpdateOnlineStorage(db *postgres.Database) *TxLoopUpdateOnlineStorage {
	return &TxLoopUpdateOnlineStorage{db: db}
}

func (s *TxLoopUpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	return s.db.WithTransaction(ctx, func(queries *dbs.Queries) error {
		for _, pair := range pairs {
			err := queries.UserOnlineUpdate(ctx, dbs.UserOnlineUpdateParams{
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
