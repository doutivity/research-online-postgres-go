package main

import (
	"context"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateOnlineStorage struct {
	repository *postgresql.Repository
}

func NewUpdateOnlineStorage(repository *postgresql.Repository) *UpdateOnlineStorage {
	return &UpdateOnlineStorage{repository: repository}
}

func (s *UpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
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
