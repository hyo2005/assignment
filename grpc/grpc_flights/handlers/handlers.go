package handlers

import (
	"assignment/grpc/grpc_flights/models"
	"assignment/grpc/grpc_flights/repositories"
	"assignment/pb"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FlightHandler struct {
	pb.UnimplementedFlightServiceServer
	flightRepository repositories.FlightRepository
}

func NewFlightHandler(flightRepo repositories.FlightRepository) (*FlightHandler, error) {
	return &FlightHandler{flightRepository: flightRepo}, nil
}

//create flight
func (h *FlightHandler) CreatFlight(ctx context.Context, in *pb.Flight) (*pb.Flight, error) {
	flight, err := h.flightRepository.CreatFlight(ctx, &models.Flight{
		ID:             uuid.New(),
		Name:           in.Name,
		From:           in.From,
		To:             in.To,
		Date:           in.GetDate().AsTime(),
		Status:         in.Status,
		Available_slot: int(in.AvailableSlot),
	})
	if err != nil {
		return nil, err
	}
	//return flight vừa đc tạo
	return &pb.Flight{
		Id:            flight.ID.String(),
		Name:          flight.Name,
		From:          flight.From,
		To:            flight.To,
		Date:          timestamppb.New(flight.Date),
		Status:        flight.Status,
		AvailableSlot: int32(flight.Available_slot),
	}, nil
}

//update flight
func (h *FlightHandler) UpdateFlight(ctx context.Context, in *pb.Flight) (*pb.Flight, error) {
	flightUpdate, err := h.flightRepository.SearchFlight(ctx, uuid.MustParse(in.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	//update thông tin cho flight
	if in.Name != "" {
		flightUpdate.Name = in.Name
	}

	if in.From != "" {
		flightUpdate.From = in.From
	}

	if in.To != "" {
		flightUpdate.To = in.To
	}

	if in.Date != nil {
		flightUpdate.Date = in.GetDate().AsTime()
	}

	if in.Status != "" {
		flightUpdate.Status = in.Status
	}

	if in.AvailableSlot != int32(flightUpdate.Available_slot) {
		flightUpdate.Available_slot = int(in.AvailableSlot)
	}

	newFlight, err := h.flightRepository.UpdateFlight(ctx, flightUpdate)
	if err != nil {
		return nil, err
	}

	return &pb.Flight{
		Id:            newFlight.ID.String(),
		Name:          newFlight.Name,
		From:          newFlight.From,
		To:            newFlight.To,
		Date:          timestamppb.New(newFlight.Date),
		Status:        newFlight.Status,
		AvailableSlot: int32(newFlight.Available_slot),
	}, nil
}

//find flight bằng id
func (h *FlightHandler) SearchFlight(ctx context.Context, in *pb.SearchFlightRequest) (*pb.Flight, error) {
	flight, err := h.flightRepository.SearchFlight(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, err
	}

	return &pb.Flight{
		Id:            flight.ID.String(),
		Name:          flight.Name,
		From:          flight.From,
		To:            flight.To,
		Date:          timestamppb.New(flight.Date),
		Status:        flight.Status,
		AvailableSlot: int32(flight.Available_slot),
	}, nil
}
