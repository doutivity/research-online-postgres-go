package postgres

import (
	"context"

	"github.com/doutivity/research-online-postgres-go/internal/storage/postgres/dbs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database - repository
type Database struct {
	pool    *pgxpool.Pool
	queries *dbs.Queries
}

// NewDatabase - constructor
func NewDatabase(pool *pgxpool.Pool) *Database {
	var queries = dbs.New(pool)

	return &Database{
		pool:    pool,
		queries: queries,
	}
}

// Connection - getter
func (r *Database) Connection() *pgxpool.Pool {
	return r.pool
}

// Queries - getter
func (r *Database) Queries() *dbs.Queries {
	return r.queries
}

// WithTransaction - start transaction
func (r *Database) WithTransaction(ctx context.Context, fn func(queries *dbs.Queries) error) error {
	return withTransaction(ctx, r.pool, r.queries, fn)
}

func withTransaction(
	ctx context.Context,
	db *pgxpool.Pool,
	queries *dbs.Queries,
	fn func(queries *dbs.Queries) error,
) (err error) {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback(ctx)

			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback(ctx)
		} else {
			// all good, commit
			err = tx.Commit(ctx)
		}
	}()

	err = fn(queries.WithTx(tx))

	return err
}
