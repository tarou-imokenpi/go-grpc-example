package main

import (
	proto "awesomeProject/pkg/proto/api"
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
)

var (
	scanner *bufio.Scanner
	client  proto.HelloServiceClient
)

func main() {
	fmt.Println("start gRPC client")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal("Connection error: ", err)
		return
	}
	defer conn.Close()

	client = proto.NewHelloServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: server stream Request")
		fmt.Println("3: exit")
		fmt.Print("Enter command number: ")

		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "1":
			sendRequest()
		case "2":
			ServerStreamRequest()
		case "3":
			fmt.Println("bye.")
			return
		}
	}
}

func ServerStreamRequest() {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	req := &proto.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("stream closed")
			break
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Response: ", res)
	}
}

func sendRequest() {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	req := &proto.HelloRequest{
		Name: name,
	}

	res, err := client.SayHello(context.Background(), req)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Response: ", res.GetMessage())
	}
}
