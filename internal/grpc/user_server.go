package grpc_server

import (
	"context"
	"go-patron/internal/usecase"
	"go-patron/proto" // Reemplaza por tu módulo real de go.mod

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server ahora tiene la dependencia del Caso de Uso inyectada
type Server struct {
	proto.UnimplementedUserServiceServer
	usecase usecase.UserUseCaseInterface
}

// NewUserGRPCServer es el constructor que inyecta la dependencia
func NewUserGRPCServer(uc *usecase.UserUseCase) *Server {
	return &Server{
		usecase: uc,
	}
}

// GetUser implementa el método RPC mapeando los datos reales de la base de datos
func (s *Server) GetUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	// 1. Extraer el ID enviado por el cliente (gRPC usualmente usa int32/int64)
	userID := uint(req.GetId())

	// 2. Llamar al Caso de Uso (Lógica de negocio y persistencia)
	// Recuerda que en pasos anteriores modificamos FindByID para que devuelva un DTO seguro (UserResponse)
	userDTO, err := s.usecase.FindByID(userID)
	if err != nil {
		// gRPC maneja códigos de error nativos (google.golang.org/grpc/status)
		// Si el error es de ID inválido, enviamos un código equivalente a Bad Request (InvalidArgument)
		if err.Error() == "ID inválido" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		// Si no se encuentra el usuario, enviamos un código NotFound de gRPC
		return nil, status.Error(codes.NotFound, "usuario no encontrado en el sistema")
	}

	// 3. Mapear el DTO obtenido del Caso de Uso hacia el UserResponse de Protobuf
	return &proto.UserResponse{
		Id:    int32(userDTO.ID),
		Name:  userDTO.Name,
		Email: userDTO.Email,
	}, nil
}
