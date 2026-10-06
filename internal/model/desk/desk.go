package desk

import (
	"lab1/internal/model"
)

func (d Desk) Status() model.BookingStatus {
	return d.status
}
