package model

import "time"

type BaseEntity struct {
	ID        int
	CreatedAt time.Time
}
