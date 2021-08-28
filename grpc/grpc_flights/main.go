package main

import (
	"assignment/grpc/grpc_flights/handlers"
	"assignment/grpc/grpc_flights/repositories"
	"assignment/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listen, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	flightRepos, err := repositories.NewDBManager()
	if err != nil {
		log.Fatal(err)
	}

	h, err := handlers.NewFlightHandler(flightRepos)
	if err != nil {
		log.Fatal(err)
	}
	reflection.Register(s)
	pb.RegisterFlightServiceServer(s, h)

	fmt.Println("Listen at port: 2222")

	s.Serve(listen)
}
