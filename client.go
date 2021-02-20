package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcChatServer/chatserver"
	// "io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)


type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)

	}
	ch.clientName = strings.TrimRight(msg, "\r\n")

}

func (ch *clientHandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		clientMessage = strings.TrimRight(clientMessage, "\r\n")
		if err != nil {
			log.Printf("Failed to read from console :: %v", err)
			continue
		}

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending to server :: %v", err)
		}

	}

}

func (ch *clientHandle) receiveMessage() {

	for {
		resp, err := ch.stream.Recv()
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		log.Printf("%s : %s", resp.Name, resp.Body)
	}
}

func main() {

	// const serverID = "34.87.108.183:5000" //"localhost:5000"


	//-------------
	fmt.Printf("Server IP:PORT ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	serverID = strings.TrimRight(serverID, "\r\n")
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}

	//-------------

	log.Println("Connecting : " + serverID)
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect gRPC server :: %v", err)
	}
	defer conn.Close()

	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to get response from gRPC server :: %v", err)
	}

	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	// block main
	bl := make(chan bool)
	<-bl

}
