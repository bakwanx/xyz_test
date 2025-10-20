package models

import "time"

type Match struct {
	ID          int64
	ScheduledAt time.Time
	HomeTeamID  int64
	AwayTeamID  int64
	HomeScore   int
	AwayScore   int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
