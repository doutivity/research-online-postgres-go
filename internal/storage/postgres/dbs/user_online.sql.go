// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: user_online.sql

package dbs

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const userOnlineUnnestUpdate = `-- name: UserOnlineUnnestUpdate :exec
UPDATE user_online AS to_t
SET online = from_t.online
FROM (
         SELECT unnest($1::BIGINT[])   AS user_id,
                unnest($2::TIMESTAMP[]) AS online
     ) AS from_t
WHERE to_t.user_id = from_t.user_id
`

type UserOnlineUnnestUpdateParams struct {
	UserIds []int64
	Onlines []pgtype.Timestamp
}

func (q *Queries) UserOnlineUnnestUpdate(ctx context.Context, arg UserOnlineUnnestUpdateParams) error {
	_, err := q.db.Exec(ctx, userOnlineUnnestUpdate, arg.UserIds, arg.Onlines)
	return err
}

const userOnlineUnnestUpsert = `-- name: UserOnlineUnnestUpsert :exec
INSERT INTO user_online (user_id, online)
VALUES (unnest($1::BIGINT[]),
        unnest($2::TIMESTAMP[]))
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online
`

type UserOnlineUnnestUpsertParams struct {
	UserIds []int64
	Onlines []pgtype.Timestamp
}

func (q *Queries) UserOnlineUnnestUpsert(ctx context.Context, arg UserOnlineUnnestUpsertParams) error {
	_, err := q.db.Exec(ctx, userOnlineUnnestUpsert, arg.UserIds, arg.Onlines)
	return err
}

const userOnlineUpdate = `-- name: UserOnlineUpdate :exec
UPDATE user_online
SET online = $1
WHERE user_id = $2
`

type UserOnlineUpdateParams struct {
	Online pgtype.Timestamp
	UserID int64
}

func (q *Queries) UserOnlineUpdate(ctx context.Context, arg UserOnlineUpdateParams) error {
	_, err := q.db.Exec(ctx, userOnlineUpdate, arg.Online, arg.UserID)
	return err
}

const userOnlineUpsert = `-- name: UserOnlineUpsert :exec
INSERT INTO user_online (user_id, online)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online
`

type UserOnlineUpsertParams struct {
	UserID int64
	Online pgtype.Timestamp
}

func (q *Queries) UserOnlineUpsert(ctx context.Context, arg UserOnlineUpsertParams) error {
	_, err := q.db.Exec(ctx, userOnlineUpsert, arg.UserID, arg.Online)
	return err
}
