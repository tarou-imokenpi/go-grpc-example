package main

import (
	hellopb "awesomeProject/pkg/proto/api"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type HelloServiceServer struct {
	hellopb.UnimplementedHelloServiceServer
}

func (s *HelloServiceServer) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	return &hellopb.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func (s *HelloServiceServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.HelloService_HelloServerStreamServer) error {
	resCount := 10
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloReply{
			Message: fmt.Sprintf("Hello %s %d", req.GetName(), i),
		}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}
	return nil
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

	hellopb.RegisterHelloServiceServer(s, NewHelloServiceServer())

	reflection.Register(s)

	go func() {
		log.Printf("Server started on port %d", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server shutting down")
	s.GracefulStop()
}
