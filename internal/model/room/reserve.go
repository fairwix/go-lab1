package room

import (
	model2 "lab1/internal/model"
)

// Reserve у Room логически отличается от Desk.Reserve: переговорку нельзя
// забронировать под группу, которая не помещается в комнату.
// Изменяет состояние -> pointer receiver.
func (r *Room) Reserve(attendees int) error {
	if r.status == model2.Booked {
		return model2.ErrAlreadyBooked
	}

	if attendees > r.Capacity {
		return model2.ValidationError{
			Message: "attendees count exceeds room capacity",
		}
	}

	r.status = model2.Booked

	return nil
}
