package main

import (
	"log"
	"os"
	"time"
	"io/ioutil"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	chat "eundoosong/grpc_chat/client/gen"
)

const (
	address     = "localhost:50051"
	defaultText = "world"
	file        = "content.txt"
)

func sendFile(client chat.ChatServiceClient) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.SendFile(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	stream.Send(&chat.File{Name: file, Type: "plain/text", Data: data})
	res, err := stream.CloseAndRecv();
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("url : %v", res)
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewChatServiceClient(conn)

	// Contact the server and print out its response.
	text := defaultText
	if len(os.Args) > 1 {
		text = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendMessage(ctx, &chat.Message{Id: "song", Text: text})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Response Id: %s", r.Id)
	log.Printf("Response Text: %s", r.Text)

	sendFile(c)
}
