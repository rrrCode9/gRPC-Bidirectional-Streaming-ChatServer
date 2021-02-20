package main

import (
	"google.golang.org/grpc"
	"grpcChatServer/chatserver"
	"log"
	"net"
	"os"
)

func main() {

	// assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "5000" // default : 5000 if in env port is not set
	}

	// initiate listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ " + Port)

	// start gRPC server instance
	grpcServer := grpc.NewServer()

	// register ChatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpcServer, &cs)

	// gRPC listen and serve
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}

}
