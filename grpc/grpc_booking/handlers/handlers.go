package handlers

import (
	"assignment/grpc/grpc_booking/models"
	"assignment/grpc/grpc_booking/repositories"
	"assignment/pb"
	"context"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingHandler struct {
	customerClient pb.CustomerServiceClient
	flightClient   pb.FlightServiceClient

	pb.UnimplementedBookingServiceServer
	bookingRepository repositories.BookingRepository
}

func NewBookingHandler(bookingRepo repositories.BookingRepository,
	customerClient pb.CustomerServiceClient,
	flightClient pb.FlightServiceClient) (*BookingHandler, error) {
	return &BookingHandler{
		bookingRepository: bookingRepo,
		customerClient:    customerClient,
		flightClient:      flightClient,
	}, nil
}

func (h *BookingHandler) CreatBooking(ctx context.Context, in *pb.Booking) (*pb.Booking, error) {
	booking, err := h.bookingRepository.CreatBooking(ctx, &models.Booking{
		ID:          uuid.New(),
		Customer_id: uuid.MustParse(in.CustomerId),
		Flight_id:   uuid.MustParse(in.FlightId),
		Code:        xid.New().String(),
		Status:      in.Status,
		Booked_date: in.GetBookedDate().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	//return flight đã được tạo
	return &pb.Booking{
		Id:         booking.ID.String(),
		CustomerId: booking.Customer_id.String(),
		FlightId:   booking.Flight_id.String(),
		Code:       booking.Code,
		Status:     booking.Status,
		BookedDate: timestamppb.New(booking.Booked_date),
	}, nil
}
func (h *BookingHandler) CancelBooking(ctx context.Context, in *pb.CancelBookingRequest) (*pb.Booking, error) {
	//tìm booking để cancel
	bookingCancel, err := h.bookingRepository.FindBooking(ctx, in.Code)
	if err != nil {
		return nil, err
	}

	//đổi trạng thái sang cancel
	if in.Status != "" {
		bookingCancel.Status = in.Status
	}
	newBooking, err := h.bookingRepository.CancelBooking(ctx, bookingCancel)
	if err != nil {
		return nil, err
	}
	//return booking vừa cancel
	return &pb.Booking{
		Id:         newBooking.ID.String(),
		CustomerId: newBooking.Customer_id.String(),
		FlightId:   newBooking.Flight_id.String(),
		Code:       newBooking.Code,
		Status:     newBooking.Status,
		BookedDate: timestamppb.New(newBooking.Booked_date),
	}, nil
}
func (h *BookingHandler) SearchBooking(ctx context.Context, in *pb.SearchBookingRequest) (*pb.BookingInfor, error) {
	booking, err := h.bookingRepository.FindBooking(ctx, in.Code)
	if err != nil {
		return nil, err
	}

	//gọi đến customer gRPC để tìm kiếm thông tin customer
	customer, err := h.customerClient.FindCustomer(ctx, &pb.FindCustomerRequest{
		Id: booking.Customer_id.String(),
	})
	if err != nil {
		return nil, err
	}

	//gọi đến flight gRPC để tìm kiếm thông tin flight
	flight, err := h.flightClient.SearchFlight(ctx, &pb.SearchFlightRequest{
		Id: booking.Flight_id.String(),
	})
	if err != nil {
		return nil, err
	}

	//return booking infor bao gồm thông tin của customer và thông tin của flight
	return &pb.BookingInfor{
		BookingDetail: &pb.Booking{
			Id:         booking.ID.String(),
			CustomerId: booking.Customer_id.String(),
			FlightId:   booking.Flight_id.String(),
			Code:       booking.Code,
			Status:     booking.Status,
			BookedDate: timestamppb.New(booking.Booked_date),
		},
		FlightDetail:   flight,
		CustomerDetail: customer,
	}, nil
}

//tìm tất cả booking của một customer
func (h *BookingHandler) SearchBookingId(ctx context.Context, in *pb.SearchBookingByIdRequest) (*pb.ListBooking, error) {
	booking, err := h.bookingRepository.FindBookingById(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, err
	}
	var pResponse pb.ListBooking
	for _, i := range *booking {
		pResponse.BookingList = append(pResponse.BookingList, &pb.Booking{
			Id:         i.ID.String(),
			CustomerId: i.Customer_id.String(),
			FlightId:   i.Flight_id.String(),
			Code:       i.Code,
			Status:     i.Status,
			BookedDate: timestamppb.New(i.Booked_date),
		})
	}
	return &pResponse, nil
}
