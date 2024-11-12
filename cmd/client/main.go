package main

import (
	hellopb "awesomeProject/pkg/proto/api"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

var (
	scanner *bufio.Scanner
	client  hellopb.HelloServiceClient
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

	client = hellopb.NewHelloServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: exit")
		fmt.Print("Enter command number: ")

		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "1":
			sendRequest()
		case "2":
			fmt.Println("bye.")
			return
		}
	}
}

func sendRequest() {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	res, err := client.SayHello(context.Background(), req)
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Response: ", res.GetMessage())
	}
}
