package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	grpcChat "github.com/kbekj/DSYS_ChittyChat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "4500", "Tcp server")

var client grpcChat.ServicesClient
var ServerConn *grpc.ClientConn

type clientHandle struct {
	stream     grpcChat.Services_ChatServiceClient
	clientName string
}

var lamportTimestamp int64

func main() {
	flag.Parse()

	file, err := os.OpenFile(fmt.Sprintf("client_%s", *clientsName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)

	fmt.Println("--- CLIENT APP ---")
	connectToServer()
	defer ServerConn.Close()

	shutDown := make(chan bool)
	msgStream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService: %v", err)
	}
	clientHandle := clientHandle{
		stream:     msgStream,
		clientName: *clientsName,
	}
	clientHandle.sendMessage("greetingSecretCode")
	go clientHandle.bindStdinToServerStream(shutDown)
	go clientHandle.streamListener(shutDown)

	<-shutDown
}

func connectToServer() {
	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", *serverPort), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	// makes a client from the server connection and saves the connection
	// and prints whether or not the connection is READY
	client = grpcChat.NewServicesClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
}

func (clientHandle *clientHandle) bindStdinToServerStream(shutDown chan bool) {

	log.Println("Binding stdin to serverstream")
	reader := bufio.NewReader(os.Stdin)

	getStdIn := func() string {
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("couldn't read from console: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")
		return clientMessage
	}

	for {
		clientMessage := getStdIn()
		for len([]rune(clientMessage)) > 128 {
			log.Printf("Too many characters current: %d Max: 128\n", len([]rune(clientMessage)))
			clientMessage = getStdIn()
		}

		clientHandle.sendMessage(clientMessage)
		if clientMessage == "bye" {
			log.Println("Shutting down")
			shutDown <- true
		}
	}
}

func (clientHandle *clientHandle) streamListener(shutDown chan bool) {
	for {
		msg, err := clientHandle.stream.Recv()
		if err != nil {
			log.Printf("Couldn't receive message from server: %v", err)
		}
		if lamportTimestamp < msg.LamportTime {
			lamportTimestamp = msg.LamportTime
		}
		log.Printf("%s : %s : at lamporttime %d \n", msg.SenderID, msg.Message, msg.LamportTime)
	}
}

func (clientHandle *clientHandle) sendMessage(clientMessage string) {

	lamportTimestamp++

	message := &grpcChat.ClientMessage{
		SenderID:    *clientsName,
		Message:     clientMessage,
		LamportTime: lamportTimestamp,
	}

	err := clientHandle.stream.Send(message)
	if err != nil {
		log.Printf("Error when sending message to stream: %v", err)
	}

}
