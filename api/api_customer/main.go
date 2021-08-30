package main

import (
	"assignment/api/api_customer/handlers"
	"assignment/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	//tạo kết nối đến grpc client
	customerConn, err := grpc.Dial(":2221", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	//Singleton pattern
	peopleClient := pb.NewCustomerServiceClient(customerConn)

	//Handler cho gin
	h := handlers.NewCustomerHandler(peopleClient)
	g := gin.Default()

	//tạo route
	g.POST("/create", h.CreatCustomer)
	g.POST("/update", h.UpdateCustomer)
	g.POST("/changepass", h.ChangePassword)
	g.GET("/find", h.FindCustomer)
	//Listen and serve
	g.Run(":4321")
}
