package desk

import (
	model2 "lab1/internal/model"
	"time"
)

type Desk struct {
	model2.BaseEntity

	Number int
	status model2.BookingStatus
}

func NewDesk(id int, number int) (Desk, error) {
	if number <= 0 {
		return Desk{}, model2.ValidationError{
			Message: "desk number must be positive",
		}
	}

	return Desk{
		BaseEntity: model2.BaseEntity{
			ID:        id,
			CreatedAt: time.Now(),
		},
		Number: number,
		status: model2.Available,
	}, nil
}
