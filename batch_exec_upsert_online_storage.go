package main

import (
	"context"
	"time"

	postgresql "github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5/pgtype"
)

type BatchExecUpsertOnlineStorage struct {
	repository *postgresql.Repository
}

func NewBatchExecUpsertOnlineStorage(repository *postgresql.Repository) *BatchExecUpsertOnlineStorage {
	return &BatchExecUpsertOnlineStorage{repository: repository}
}

func (s *BatchExecUpsertOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	args := make([]dbs.UserOnlineBatchExecUpsertParams, len(pairs))
	for i, pair := range pairs {
		args[i] = dbs.UserOnlineBatchExecUpsertParams{
			Online: pgtype.Timestamp{
				Time:  time.Unix(pair.Timestamp, 0).UTC(),
				Valid: true,
			},
			UserID: pair.UserID,
		}
	}

	var batchErr error

	s.repository.Queries().UserOnlineBatchExecUpsert(ctx, args).Exec(func(i int, err error) {
		if err != nil {
			batchErr = multierror.Append(batchErr, err)
		}
	})

	return batchErr
}
