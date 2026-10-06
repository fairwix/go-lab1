package desk

import (
	model2 "lab1/internal/model"
)

// У Desk параметр attendees осознанно не используется: стол один и тот же
// вне зависимости от количества людей — в этом и есть логическое отличие
// от Room.Reserve, где attendees определяет, пройдёт ли валидация.
func (d *Desk) Reserve(_ int) error {
	if d.status == model2.Booked {
		return model2.ErrAlreadyBooked
	}

	d.status = model2.Booked

	return nil
}
