package room

import "lab1/internal/model"

func (r Room) Status() model.BookingStatus {
	return r.status
}
