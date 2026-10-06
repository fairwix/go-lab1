package model

type BookingStatus int

const (
	Available BookingStatus = iota
	Booked
)

func (s BookingStatus) String() string {
	switch s {
	case Available:
		return "available"
	case Booked:
		return "booked"
	default:
		return "unknown"
	}
}
