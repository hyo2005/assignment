package main

import (
	"assignment/grpc/grpc_customer/handlers"
	"assignment/grpc/grpc_customer/repositories"
	"assignment/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	bookingConnect, err := grpc.Dial(":2223", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	bookingClient := pb.NewBookingServiceClient(bookingConnect)

	listen, err := net.Listen("tcp", ":2221")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	customerRepos, err := repositories.NewDBManager()
	if err != nil {
		log.Fatal(err)
	}

	h, err := handlers.NewCustomerHandler(customerRepos, bookingClient)
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(s)
	pb.RegisterCustomerServiceServer(s, h)

	fmt.Println("Listen at port: 2221")

	s.Serve(listen)
}
