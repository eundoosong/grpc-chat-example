package main

import (
	"log"
	"time"
	"io/ioutil"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	chat "eundoosong/grpc-examples/client/go/gen"
	"fmt"
	"os"
)

const (
	address     = "localhost:50051"
	defaultUser = "song"
	file        = "content.txt"
)

func menu() {
	fmt.Println("1. send a message")
	fmt.Println("2. send a file")
	fmt.Println("99. exit")
}

func sendMessage(client chat.ChatServiceClient, user string, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.SendMessage(ctx, &chat.Message{Id: user, Text: text})
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	log.Printf("Response User: %s", r.Id)
	log.Printf("Response Text: %s", r.Text)

}

func sendFile(client chat.ChatServiceClient) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.UploadFiles(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
  stream.Send(&chat.File{Name: file, Type: "plain/text", Len: int32(len(data)), Data: data})
	res, err := stream.CloseAndRecv();
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	log.Printf("url : %v", res)
}

func main() {
	user := defaultUser
	if len(os.Args) > 1 {
		user = os.Args[1]
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := chat.NewChatServiceClient(conn)

	var input string
	for {
		menu();
		fmt.Scanln(&input)
		if input == "1" {
			var text string;
			fmt.Print("message("+user+"): ")
			fmt.Scanln(&text)
			sendMessage(client, user, text)
		} else if input == "2" {
			sendFile(client)
		} else {
			break;
		}

	}
}
