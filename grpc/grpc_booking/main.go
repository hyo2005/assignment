package main

import (
	"assignment/grpc/grpc_booking/handlers"
	"assignment/grpc/grpc_booking/repositories"
	"assignment/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	customerConnect, err := grpc.Dial(":2221", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	flightConnect, err := grpc.Dial(":2222", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	customerClient := pb.NewCustomerServiceClient(customerConnect)
	flightClient := pb.NewFlightServiceClient(flightConnect)

	listen, err := net.Listen("tcp", ":2223")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	bookingRepos, err := repositories.NewDBManager()
	if err != nil {
		log.Fatal(err)
	}

	h, err := handlers.NewBookingHandler(bookingRepos, customerClient, flightClient)
	if err != nil {
		log.Fatal(err)
	}
	reflection.Register(s)
	pb.RegisterBookingServiceServer(s, h)

	fmt.Println("Listen at port: 2223")

	s.Serve(listen)
}
