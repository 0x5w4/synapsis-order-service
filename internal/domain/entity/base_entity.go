package entity

import "time"

type Base struct {
	ID        uint32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
