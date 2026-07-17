package main

import (
	"context"
	"fmt"
	"net"

	"go-patron/proto"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedUserServiceServer
}

// Método
func (s *Server) GetUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {

	return &proto.UserResponse{
		Id:    req.Id,
		Name:  "Juan",
		Email: "juan@gmail.com",
	}, nil
}

func main() {

	lis, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	proto.RegisterUserServiceServer(
		server,
		&Server{},
	)
	fmt.Println("Servidor corriendo !")
	server.Serve(lis)
}
