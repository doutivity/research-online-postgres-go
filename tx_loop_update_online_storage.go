package main

import (
	"context"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type TxLoopUpdateOnlineStorage struct {
	repository *postgresql.Repository
}

func NewTxLoopUpdateOnlineStorage(repository *postgresql.Repository) *TxLoopUpdateOnlineStorage {
	return &TxLoopUpdateOnlineStorage{repository: repository}
}

func (s *TxLoopUpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	return s.repository.WithTransaction(ctx, func(queries *dbs.Queries) error {
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
