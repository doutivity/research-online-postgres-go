package main

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
