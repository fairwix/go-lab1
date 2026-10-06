package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	model "lab1/internal/model"
	"lab1/internal/service"
)

type Handler struct {
	service service.BookingService
}

func NewHandler(service service.BookingService) *Handler {
	return &Handler{
		service: service,
	}
}

type BookRequest struct {
	ID        int `json:"id"`
	Attendees int `json:"attendees"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	states := h.service.GetState()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(states); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}

func (h *Handler) Book(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var request BookRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(
			w,
			"invalid JSON",
			http.StatusBadRequest,
		)
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(
			w,
			"invalid JSON",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Book(request.ID, request.Attendees); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {
	var validationErr model.ValidationError

	switch {
	case errors.As(err, &validationErr):
		http.Error(w, validationErr.Error(), http.StatusBadRequest)

	case errors.Is(err, model.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, model.ErrAlreadyBooked):
		http.Error(w, err.Error(), http.StatusConflict)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
