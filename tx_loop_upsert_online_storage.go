package main

import (
	"context"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type TxLoopUpsertOnlineStorage struct {
	repository *postgresql.Repository
}

func NewTxLoopUpsertOnlineStorage(repository *postgresql.Repository) *TxLoopUpsertOnlineStorage {
	return &TxLoopUpsertOnlineStorage{repository: repository}
}

func (s *TxLoopUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	return s.repository.WithTransaction(ctx, func(queries *dbs.Queries) error {
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
