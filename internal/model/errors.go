package model

import "errors"

var ErrAlreadyBooked = errors.New("resource is already booked")
var ErrNotFound = errors.New("resource not found")

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
