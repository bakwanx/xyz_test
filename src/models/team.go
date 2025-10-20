package models

import "time"

type Team struct {
	ID          int64
	Name        string
	LogoPath    *string
	FoundedYear int
	Address     string
	City        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
