package main

import (
	"awesomeProject/pkg/proto"
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
	proto.UnimplementedHelloServiceServer
}

func (s *HelloServiceServer) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

func (s *HelloServiceServer) HelloServerStream(req *proto.HelloRequest, stream proto.HelloService_HelloServerStreamServer) error {
	resCount := 10
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&proto.HelloReply{
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

	proto.RegisterHelloServiceServer(s, NewHelloServiceServer())

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
