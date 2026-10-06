package service

import (
	"fmt"
	"sort"
	"sync"

	model "lab1/internal/model"
)

type BookingService interface {
	Book(id int, attendees int) error
	GetState() []ResourceState
}

type ResourceState struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type bookingService struct {
	mu        sync.Mutex
	resources map[int]model.Reservable
}

var _ BookingService = (*bookingService)(nil)

func NewBookingService(resources map[int]model.Reservable) BookingService {
	return &bookingService{
		resources: resources,
	}
}

func (s *bookingService) Book(id int, attendees int) error {
	// Это бизнес-правило, одинаковое для HTTP, gRPC и любого будущего транспорта.
	if id <= 0 {
		return model.ValidationError{
			Message: "resource id must be positive",
		}
	}
	if attendees <= 0 {
		return model.ValidationError{
			Message: "attendees count must be positive",
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	resource, ok := s.resources[id]
	if !ok {
		return model.ErrNotFound
	}

	if err := resource.Reserve(attendees); err != nil {
		return fmt.Errorf("failed to book resource %d: %w", id, err)
	}

	return nil
}

func (s *bookingService) GetState() []ResourceState {
	s.mu.Lock()
	defer s.mu.Unlock()

	states := make([]ResourceState, 0, len(s.resources))

	for id, resource := range s.resources {
		states = append(states, ResourceState{
			ID:     id,
			Status: resource.Status().String(),
		})
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].ID < states[j].ID
	})

	return states
}
