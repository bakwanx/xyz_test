package models

import "time"

type Player struct {
	ID           int64
	TeamID       int64
	Name         string
	HeightCm     int
	WeightKg     int
	Position     string
	JerseyNumber int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
