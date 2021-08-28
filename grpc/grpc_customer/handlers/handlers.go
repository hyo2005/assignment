package handlers

import (
	"assignment/grpc/grpc_customer/models"
	"assignment/grpc/grpc_customer/repositories"
	"assignment/pb"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerHandler struct {
	bookingClient pb.BookingServiceClient
	pb.UnimplementedCustomerServiceServer
	customerRepository repositories.CustomerRepository
}

func NewCustomerHandler(customerRepo repositories.CustomerRepository,
	bookingClient pb.BookingServiceClient) (*CustomerHandler, error) {
	return &CustomerHandler{
		customerRepository: customerRepo,
		bookingClient:      bookingClient,
	}, nil
}
func (h *CustomerHandler) CreateCustomer(ctx context.Context, in *pb.Customer) (*pb.Customer, error) {
	//tạo customer mới
	customer, err := h.customerRepository.CreateCustomer(ctx, &models.Customer{
		ID:         uuid.New(),
		Name:       in.CustomerName,
		Address:    in.Address,
		Phone:      in.Phone,
		License_id: in.LicenseId,
		Email:      in.Email,
		Password:   in.Password,
		Active:     in.Active,
	})
	if err != nil {
		return nil, err
	}
	//return customer vừa đc tạo
	return &pb.Customer{
		Id:           customer.ID.String(),
		CustomerName: customer.Name,
		Address:      customer.Address,
		Phone:        customer.Phone,
		LicenseId:    customer.License_id,
		Email:        customer.Email,
		Password:     customer.Password,
		Active:       customer.Active,
	}, nil
}

//update customer
func (h *CustomerHandler) UpdateCustomer(ctx context.Context, in *pb.Customer) (*pb.Customer, error) {
	customerUpdated, err := h.customerRepository.FindCustomer(ctx, uuid.MustParse(in.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	//update thông tin customer
	if in.CustomerName != "" {
		customerUpdated.Name = in.CustomerName
	}

	if in.Address != "" {
		customerUpdated.Address = in.Address
	}

	if in.Phone != "" {
		customerUpdated.Phone = in.Phone
	}

	if in.LicenseId != "" {
		customerUpdated.License_id = in.LicenseId
	}

	if in.Email != "" {
		customerUpdated.Email = in.Email
	}

	if in.Active != customerUpdated.Active {
		customerUpdated.Active = in.Active
	}

	//update tới db
	newCustomer, err := h.customerRepository.UpdateCustomer(ctx, customerUpdated)
	if err != nil {
		return nil, err
	}

	return &pb.Customer{
		CustomerName: newCustomer.Name,
		Address:      newCustomer.Address,
		Phone:        newCustomer.Phone,
		LicenseId:    newCustomer.License_id,
		Email:        newCustomer.Email,
		Active:       newCustomer.Active,
		Id:           newCustomer.ID.String(),
	}, nil
}

//Thay đổi Password
func (h *CustomerHandler) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Customer, error) {
	customerUpdated, err := h.customerRepository.FindCustomer(ctx, uuid.MustParse(in.Id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	//update customer password
	customerUpdated.Password = in.Password

	//update tới db
	newCustomer, err := h.customerRepository.UpdateCustomer(ctx, customerUpdated)
	if err != nil {
		return nil, err
	}

	return &pb.Customer{
		CustomerName: newCustomer.Name,
		Address:      newCustomer.Address,
		Phone:        newCustomer.Phone,
		LicenseId:    newCustomer.License_id,
		Active:       newCustomer.Active,
		Id:           newCustomer.ID.String(),
	}, nil
}

func (h *CustomerHandler) FindCustomer(ctx context.Context, in *pb.FindCustomerRequest) (*pb.Customer, error) {
	customer, err := h.customerRepository.FindCustomer(ctx, uuid.MustParse(in.Id))
	if err != nil {
		return nil, err
	}

	//return customer vừa tìm được
	return &pb.Customer{
		CustomerName: customer.Name,
		Address:      customer.Address,
		Phone:        customer.Phone,
		LicenseId:    customer.License_id,
		Active:       customer.Active,
		Id:           customer.ID.String(),
		Email:        customer.Email,
		Password:     customer.Password,
	}, nil
}
func (h *CustomerHandler) BookingHistory(ctx context.Context, in *pb.BookingHistoryRequest) (*pb.BookingHistoryResponse, error) {
	bookings, err := h.bookingClient.SearchBookingId(ctx, &pb.SearchBookingByIdRequest{
		Id: in.Id,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	var pResponse pb.BookingHistoryResponse
	pResponse.CustomerId = in.Id
	for _, v := range bookings.BookingList {
		pResponse.Code = append(pResponse.Code, v.Code)
		pResponse.Id = append(pResponse.Id, v.Id)
		pResponse.FlightId = append(pResponse.FlightId, v.FlightId)
		pResponse.Status = append(pResponse.Status, v.Status)
		pResponse.BookedDate = append(pResponse.BookedDate, v.BookedDate)
	}

	return &pResponse, nil
}
