package main

import (
	"io/ioutil"
	"log"
	"time"
	"io"

	chat "eundoosong/grpc-examples/client/go/gen"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
)

const (
	address     = "localhost:50051"
	defaultUser = "song"
	file        = "content.txt"
)

func menu() {
	fmt.Println("1. send a message")
	fmt.Println("2. upload files")
	fmt.Println("3. download files")
	fmt.Println("4. transcode files")
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

func uploadFiles(client chat.ChatServiceClient) {
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
	for i := 0; i < 10; i++ {
		if err := stream.Send(&chat.File{Name: file,
													 Type: "plain/text",
													 Len: int32(len(data)),
													 Data: data}); err != nil {
			log.Fatalf("%v.Send = %v", stream, err)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("file length : %d", len(res.Id))
	for id := range res.Id {
		log.Printf("id: %#s", id)
	}
}

func downloadFiles(client chat.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	list := []string{"aaaa", "bbbb"}
	stream, err := client.DownloadFiles(ctx, &chat.FileIds{Id: list})
	if err != nil {
		log.Fatalf("%v.downloadFiles(_) = _, %v", client, err)
	}

	for {
		file, err := stream.Recv()
		if err == io.EOF {
			log.Printf("EOF")
			break
		}
		if err != nil {
			log.Fatalf("%v.downloadFiles(_) = _, %v", client, err)
		}
		log.Println(file.Name)
	}
}

func TranscodeFiles(client chat.ChatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.ConvertFiles(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	waitch := make(chan struct{})
	go func() {
		for {
			file, err := stream.Recv()
			if err == io.EOF {
				log.Printf("EOF")
				close(waitch)
				break
			}
			if err != nil {
				log.Fatalf("%v.downloadFiles(_) = _, %v", client, err)
			}
			log.Println(file.Name)
		}
	}()

	for i := 0; i < 10; i++ {
		if err := stream.Send(&chat.File{Name: file,
													 Type: "plain/text",
													 Len: int32(len(data)),
													 Data: data}); err != nil {
			log.Fatalf("%v.Send = %v", stream, err)
		}
	}
	stream.CloseSend()
	<-waitch
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
		menu()
		fmt.Scanln(&input)
		if input == "1" {
			var text string
			fmt.Print("message(" + user + "): ")
			fmt.Scanln(&text)
			sendMessage(client, user, text)
		} else if input == "2" {
			uploadFiles(client)
		} else if input == "3" {
			downloadFiles(client)
		} else if input == "4" {
			TranscodeFiles(client)
		} else {
			break
		}
	}
}
