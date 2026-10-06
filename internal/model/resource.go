// internal/model/resource.go
package model

// Reserve принимает attendees — для Room это осмысленный параметр
// (сколько человек придёт), для Desk он игнорируется (стол — на одного).
type Reservable interface {
	Reserve(attendees int) error
	Status() BookingStatus
}
