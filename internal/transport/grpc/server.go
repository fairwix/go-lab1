package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	model "lab1/internal/model"
	"lab1/internal/service"
	pb "lab1/proto"
)

type Server struct {
	pb.UnimplementedBookingServiceServer

	service service.BookingService
}

func NewServer(service service.BookingService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Book(
	ctx context.Context,
	req *pb.BookRequest,
) (*pb.BookResponse, error) {

	if err := s.service.Book(int(req.GetId()), int(req.GetAttendees())); err != nil {
		return nil, mapError(err)
	}

	return &pb.BookResponse{}, nil
}

func (s *Server) GetState(
	ctx context.Context,
	req *pb.GetStateRequest,
) (*pb.GetStateResponse, error) {

	states := s.service.GetState()

	resources := make(
		[]*pb.ResourceState,
		0,
		len(states),
	)

	for _, state := range states {
		resources = append(resources, &pb.ResourceState{
			Id:     int32(state.ID),
			Status: state.Status,
		})
	}

	return &pb.GetStateResponse{
		Resources: resources,
	}, nil
}

func mapError(err error) error {
	var validationErr model.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return status.Error(
			codes.InvalidArgument,
			validationErr.Error(),
		)

	case errors.Is(err, model.ErrNotFound):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case errors.Is(err, model.ErrAlreadyBooked):
		return status.Error(
			codes.AlreadyExists,
			err.Error(),
		)

	default:
		return status.Error(
			codes.Internal,
			"internal server error",
		)
	}
}
