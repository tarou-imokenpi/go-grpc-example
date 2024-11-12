package main

import (
	hello "awesomeProject/pkg/proto/api"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type HelloServiceServer struct {
	hello.UnimplementedHelloServiceServer
}

func (s *HelloServiceServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	return &hello.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func (s *HelloServiceServer) HelloServerStream(req *hello.HelloRequest, stream hello.HelloService_HelloServerStreamServer) error {
	resCount := 10
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hello.HelloReply{
			Message: fmt.Sprintf("Hello %s %d", req.GetName(), i),
		}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}

	// ストリームの終了
	return nil
}

func (s *HelloServiceServer) HelloClientStream(stream hello.HelloService_HelloClientStreamServer) error {
	nameList := make([]string, 0)

	for {
		req, err := stream.Recv()

		// クライアントからのリクエストをすべて受け取ったとき
		// クライアントに名前のリストを返す
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello %v", nameList)
			return stream.SendAndClose(&hello.HelloReply{
				Message: message,
			})
		}
		if err != nil {
			return err
		}

		nameList = append(nameList, req.GetName())
	}
}

func NewHelloServiceServer() *HelloServiceServer {
	return &HelloServiceServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	hello.RegisterHelloServiceServer(s, NewHelloServiceServer())

	reflection.Register(s)

	go func() {
		log.Printf("Server started on port %d", port)
		err := s.Serve(listener)
		if err != nil {
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server shutting down")
	s.GracefulStop()
}
