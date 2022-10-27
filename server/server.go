package main

import (
    "io"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	grpcChat "github.com/kbekj/DSYS_ChittyChat/proto"
	"google.golang.org/grpc"
)

type message struct {
	MessageBody string
	SenderID    string
}

type messages struct {
	messageQue []message
	mutex      sync.Mutex
}

type Server struct {
	grpcChat.UnimplementedServicesServer // an interface that the server needs to have

	name string
	port string
}

type connectedClient struct {
	name   string
	stream grpcChat.Services_ChatServiceServer
}

// TODO: Should this be stored in server
var messagesObject = messages{}
var connectedClientStreams = []connectedClient{}
var connectedClientsAmount int
var lamportTimestamp int64

var serverName = flag.String("name", "default", "Senders name")
var port = flag.String("port", "4500", "Server port")

func main() {
	file, err := os.OpenFile(fmt.Sprintf("server_%s", *serverName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	flag.Parse()
	log.Println(".:server is starting:.")
	launchServer()
}

func launchServer() {
	connectedClientsAmount = 0
	lamportTimestamp = 0
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{
		name: *serverName,
		port: *port,
	}
	grpcChat.RegisterServicesServer(grpcServer, server)
	log.Printf("Server %s: Listening at %v\n", *serverName, list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func (s *Server) ChatService(msgStream grpcChat.Services_ChatServiceServer) error {
	connectedClientID := strconv.Itoa(connectedClientsAmount)
	connectedClientStreams = append(connectedClientStreams, connectedClient{
		stream: msgStream,
		name:   connectedClientID,
	})
	connectedClientsAmount++ //TODO: not atomic with read above

	errorChannel := make(chan error)
	go receiveStream(msgStream, connectedClientID, errorChannel)
	go messagesListener(msgStream, errorChannel)

	return <-errorChannel
}

func receiveStream(msgStream grpcChat.Services_ChatServiceServer, connectedClientID string, errorChannel chan error) {
	for {
		msg, err := msgStream.Recv()
		if err != nil {
			log.Printf("Something went wrong: %v", err)
			errorChannel <- err
			return
		}
		if msg.LamportTime > lamportTimestamp {
			lamportTimestamp = msg.LamportTime
		}
		switch msg.Message {
		case "bye":
			var streamIndex int
			for i := 0; i < len(connectedClientStreams); i++ {
				if connectedClientStreams[i].name == connectedClientID {
					streamIndex = i
					break
				}
			}
			connectedClientStreams[streamIndex] = connectedClientStreams[len(connectedClientStreams)-1]
			connectedClientStreams = connectedClientStreams[:len(connectedClientStreams)-1]
			sendToAllStreams(*serverName, fmt.Sprintf("Participant %s left Chitty-Chat at Lamport time %d", msg.SenderID, lamportTimestamp))
			errorChannel <- err
            return
		case "greetingSecretCode":
			sendToAllStreams(*serverName, fmt.Sprintf("Participant %s joined Chitty-Chat at Lamport time %d", msg.SenderID, lamportTimestamp+1))
		default:
			messagesObject.mutex.Lock()

			messagesObject.messageQue = append(messagesObject.messageQue, message{
				MessageBody: msg.Message,
				SenderID:    msg.SenderID,
			})
			messagesObject.mutex.Unlock()
			objectBodyReceived := messagesObject.messageQue[len(messagesObject.messageQue)-1]
			log.Printf("Message recieved as: %s\nfrom: %s\n at lamport time stamp: %d\n", objectBodyReceived.MessageBody, objectBodyReceived.SenderID, lamportTimestamp)
		}
	}
}

func messagesListener(msgStream grpcChat.Services_ChatServiceServer, errorChannel chan error) {
	for {
		time.Sleep(500 * time.Millisecond)

		messagesObject.mutex.Lock()

		if len(messagesObject.messageQue) == 0 {
			messagesObject.mutex.Unlock()
			continue
		}

		senderID := messagesObject.messageQue[0].SenderID
		newMessage := messagesObject.messageQue[0].MessageBody
		messagesObject.mutex.Unlock()

		err := sendToAllStreams(senderID, newMessage)
		if err != nil {
			errorChannel <- err
		}

		messagesObject.mutex.Lock()
		if len(messagesObject.messageQue) > 1 {
			messagesObject.messageQue = messagesObject.messageQue[1:]
		} else {
			messagesObject.messageQue = []message{}
		}
		messagesObject.mutex.Unlock()
	}
}

func sendToAllStreams(senderID string, newMessage string) error {
	lamportTimestamp++
	for _, v := range connectedClientStreams {
		err := v.stream.Send(&grpcChat.ServerMessage{
			SenderID:    senderID,
			Message:     newMessage,
			LamportTime: lamportTimestamp,
		})
		if err != nil {
			return err
		}

	}
	return nil
}
