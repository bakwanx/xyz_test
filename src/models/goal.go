package models

import "time"

type Goal struct {
	ID           int64
	MatchID      int64
	PlayerID     int64
	TeamID       int64
	MinuteScored int
	CreatedAt    time.Time
	Player       struct {
		Name         string
		JerseyNumber int
	}
}
