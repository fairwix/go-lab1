package room

import (
	model2 "lab1/internal/model"
	"time"
)

type Room struct {
	model2.BaseEntity

	Name     string
	Capacity int
	status   model2.BookingStatus
}

func NewRoom(id int, name string, capacity int) (Room, error) {
	if name == "" {
		return Room{}, model2.ValidationError{
			Message: "room name cannot be empty",
		}
	}

	if capacity <= 0 {
		return Room{}, model2.ValidationError{
			Message: "room capacity must be positive",
		}
	}

	return Room{
		BaseEntity: model2.BaseEntity{
			ID:        id,
			CreatedAt: time.Now(),
		},
		Name:     name,
		Capacity: capacity,
		status:   model2.Available,
	}, nil
}
