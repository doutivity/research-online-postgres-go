// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user_online.sql

package dbs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const userOnlineNew = `-- name: UserOnlineNew :exec
INSERT INTO user_online (user_id, online)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
    SET online = $2
`

type UserOnlineNewParams struct {
	UserOnline int64
	Online     pgtype.Timestamptz
}

func (q *Queries) UserOnlineNew(ctx context.Context, arg UserOnlineNewParams) error {
	_, err := q.db.Exec(ctx, userOnlineNew, arg.UserOnline, arg.Online)
	return err
}
