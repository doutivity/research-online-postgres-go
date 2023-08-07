package main

import (
	"context"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpsertOnlineStorage struct {
	repository *postgresql.Repository
}

func NewUpsertOnlineStorage(repository *postgresql.Repository) *UpsertOnlineStorage {
	return &UpsertOnlineStorage{repository: repository}
}

func (s *UpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	return s.repository.WithTransaction(ctx, func(queries *dbs.Queries) error {
		for _, pair := range pairs {
			err := queries.UserOnlineNew(ctx, dbs.UserOnlineNewParams{
				UserID: pair.UserID,
				Online: pgtype.Timestamptz{
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
