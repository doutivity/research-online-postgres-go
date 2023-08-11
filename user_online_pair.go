package main

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type UserOnlinePair struct {
	UserID    int64
	Timestamp int64
}

func incUserOnlinePair(source UserOnlinePair, shift int64) UserOnlinePair {
	return UserOnlinePair{
		UserID:    source.UserID,
		Timestamp: source.Timestamp + shift,
	}
}

func userOnlinePairsToPgxSlices(pairs []UserOnlinePair) ([]int64, []pgtype.Timestamp) {
	var (
		userIDs    = make([]int64, len(pairs))
		timestamps = make([]pgtype.Timestamp, len(pairs))
	)

	for i, pair := range pairs {
		userIDs[i] = pair.UserID
		timestamps[i] = pgtype.Timestamp{
			Time:  time.Unix(pair.Timestamp, 0).UTC(),
			Valid: true,
		}
	}

	return userIDs, timestamps
}
