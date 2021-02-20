package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue        []messageUnit
	clientCount int
	mu          sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

// ChatService
func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {

	clientUniqueCode := rand.Intn(1e3)

	// recieve request <<< client
	go recieveFromStream(csi, clientUniqueCode)
	//stream >>> client
	errch := make(chan error)
	go sendToStream(csi, clientUniqueCode, errch)

	return <-errch
}

// recieve from stream
func recieveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int) {

	for {
		req, err := csi_.Recv()
		if err != nil {
			log.Printf("Error reciving request from client :: %v", err)
			break

		} else {
			messageHandleObject.mu.Lock()
			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{ClientName: req.Name, MessageBody: req.Body, MessageUniqueCode: rand.Intn(1e8), ClientUniqueCode: clientUniqueCode_})
			messageHandleObject.mu.Unlock()
			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])
		}

	}

}

//send to stream
func sendToStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {

	for {

		for {
			time.Sleep(500 * time.Millisecond)
			messageHandleObject.mu.Lock()
			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}
			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4client := messageHandleObject.MQue[0].ClientName
			message4client := messageHandleObject.MQue[0].MessageBody
			messageHandleObject.mu.Unlock()
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&FromServer{Name: senderName4client, Body: message4client})

				if err != nil {
					errch_ <- err
				}
				messageHandleObject.mu.Lock()
				if len(messageHandleObject.MQue) >= 2 {
					messageHandleObject.MQue = messageHandleObject.MQue[1:] // if send success > delete message
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}
				messageHandleObject.mu.Unlock()

			}

		}

		time.Sleep(100 * time.Millisecond)

	}

}
