package main

import (
	"context"
	"time"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres"
	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5/pgtype"
)

type BatchExecUpdateOnlineStorage struct {
	db *postgres.Database
}

func NewBatchExecUpdateOnlineStorage(db *postgres.Database) *BatchExecUpdateOnlineStorage {
	return &BatchExecUpdateOnlineStorage{db: db}
}

func (s *BatchExecUpdateOnlineStorage) BatchStore(ctx context.Context, pairs []UserOnlinePair) error {
	args := make([]dbs.UserOnlineBatchExecUpdateParams, len(pairs))
	for i, pair := range pairs {
		args[i] = dbs.UserOnlineBatchExecUpdateParams{
			Online: pgtype.Timestamp{
				Time:  time.Unix(pair.Timestamp, 0).UTC(),
				Valid: true,
			},
			UserID: pair.UserID,
		}
	}

	var batchErr error

	s.db.Queries().UserOnlineBatchExecUpdate(ctx, args).Exec(func(i int, err error) {
		if err != nil {
			batchErr = multierror.Append(batchErr, err)
		}
	})

	return batchErr
}
