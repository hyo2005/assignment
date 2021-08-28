package handlers

import (
	"assignment/api/api_customer/requests"
	"assignment/api/api_customer/responses"
	"assignment/pb"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerHandler interface {
	CreatCustomer(c *gin.Context)
	UpdateCustomer(c *gin.Context)
	ChangePassword(c *gin.Context)
	FindCustomer(c *gin.Context)
	BookingHistory(c *gin.Context)
}

type customerHandler struct {
	customerClient pb.CustomerServiceClient
}

func NewCustomerHandler(customerClient pb.CustomerServiceClient) CustomerHandler {
	return &customerHandler{
		customerClient: customerClient,
	}
}
func (ch *customerHandler) CreatCustomer(c *gin.Context) {
	req := requests.CreatCustomerRequest{} //khai báo CreatCustomerRequest object

	//parse form từ http request
	if err := c.ShouldBindJSON(&req); err != nil {
		//validate form
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			errMessage := make([]string, 0)
			for _, v := range validateErr {
				errMessage = append(errMessage, v.Error())
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusText(http.StatusBadRequest),
				"error":  errMessage,
			})

			return
		}
	}

	//protobuf model customer
	pRequest := &pb.Customer{
		CustomerName: req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		LicenseId:    req.License_id,
		Active:       req.Active,
		Email:        req.Email,
		Password:     req.Password,
	}

	//gRPC client gọi đến hàm CreateCustomer
	pResponse, err := ch.customerClient.CreateCustomer(c.Request.Context(), pRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusText(http.StatusInternalServerError),
			"error":  err.Error(),
		})
		return
	}

	//response object
	dto := &responses.CustomerResponse{
		Name:       pResponse.CustomerName,
		Phone:      pResponse.Phone,
		License_id: pResponse.LicenseId,
		Address:    pResponse.Address,
		Email:      pResponse.Email,
		ID:         pResponse.Id,
	}

	//json to response http
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusText(http.StatusOK),
		"payload": dto,
	})
}
func (ch *customerHandler) UpdateCustomer(c *gin.Context) {
	req := requests.UpdateCustomerRequest{}

	//bind form to req
	if err := c.ShouldBindJSON(&req); err != nil {
		//validate form
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			errMessage := make([]string, 0)
			for _, v := range validateErr {
				errMessage = append(errMessage, v.Error())
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusText(http.StatusBadRequest),
				"error":  errMessage,
			})

			return
		}
	}

	//request proto
	pRequest := pb.Customer{
		CustomerName: req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		LicenseId:    req.License_id,
		Id:           req.ID,
		Email:        req.Email,
		Active:       req.Active,
	}

	//update trên grpc
	pResponse, err := ch.customerClient.UpdateCustomer(c, &pRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusText(http.StatusInternalServerError),
			"error":  err.Error(),
		})
		return
	}

	//dto chứa kết quả trả về cho client
	dto := responses.CustomerResponse{
		ID:         pResponse.Id,
		Name:       pResponse.CustomerName,
		Phone:      pResponse.Phone,
		License_id: pResponse.LicenseId,
		Address:    pResponse.Address,
		Email:      pResponse.Email,
		Active:     pResponse.Active,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusText(http.StatusOK),
		"payload": dto,
	})
}
func (ch *customerHandler) FindCustomer(c *gin.Context) {
	req := requests.FindCustomerRequest{} //khai báo FindCustomerRequest object

	//parse form request
	if err := c.ShouldBindJSON(&req); err != nil {
		//validate form
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			errMessage := make([]string, 0)
			for _, v := range validateErr {
				errMessage = append(errMessage, v.Error())
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusText(http.StatusBadRequest),
				"error":  errMessage,
			})

			return
		}
	}

	//protobuf model find customer request
	pRequest := pb.FindCustomerRequest{
		Id: req.ID,
	}

	//gRPC client gọi hàm find customer
	pResponse, err := ch.customerClient.FindCustomer(c.Request.Context(), &pRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusText(http.StatusInternalServerError),
			"error":  err.Error(),
		})
		return
	}

	//response object
	dto := responses.CustomerResponse{
		ID:         pResponse.Id,
		Name:       pResponse.CustomerName,
		Phone:      pResponse.Phone,
		License_id: pResponse.LicenseId,
		Address:    pResponse.Address,
		Email:      pResponse.Email,
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusText(http.StatusOK),
		"payload": dto,
	})
}
func (ch *customerHandler) ChangePassword(c *gin.Context) {
	req := requests.ChangePasswordRequest{}
	//binding form
	if err := c.ShouldBindJSON(&req); err != nil {
		//validate form
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			errMessage := make([]string, 0)
			for _, v := range validateErr {
				errMessage = append(errMessage, v.Error())
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusText(http.StatusBadRequest),
				"error":  errMessage,
			})

			return
		}
	}

	//proto request
	pRequest := pb.ChangePasswordRequest{
		Id:       req.ID,
		Password: req.Password,
	}

	_, err := ch.customerClient.ChangePassword(c, &pRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusText(http.StatusInternalServerError),
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusText(http.StatusOK),
		"payload": "Change password success",
	})
}
func (ch *customerHandler) BookingHistory(c *gin.Context) {
	req := requests.FindCustomerRequest{}

	//parse form từ http request
	if err := c.ShouldBindJSON(&req); err != nil {
		//validate form
		if validateErr, ok := err.(validator.ValidationErrors); ok {
			errMessage := make([]string, 0)
			for _, v := range validateErr {
				errMessage = append(errMessage, v.Error())
			}

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusText(http.StatusBadRequest),
				"error":  errMessage,
			})

			return
		}
	}

	pCRequest := &pb.BookingHistoryRequest{
		Id: req.ID,
	}

	pResponse, err := ch.customerClient.BookingHistory(c.Request.Context(), pCRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusText(http.StatusInternalServerError),
			"error":  err.Error(),
		})
		return
	}

	dto := responses.HistoryResponse{}
	dto.Customrer_id = pResponse.CustomerId

	for _, v := range pResponse.BookedDate {
		dto.Booking_date = append(dto.Booking_date, v.AsTime())
	}
	dto.Booking_code = append(dto.Booking_code, pResponse.Code...)
	dto.Flight_id = append(dto.Flight_id, pResponse.FlightId...)
	dto.Booking_status = append(dto.Booking_status, pResponse.Status...)
	dto.Booking_id = append(dto.Booking_id, pResponse.Id...)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusText(http.StatusOK),
		"payload": dto,
	})
}
