package grpc_server

import (
	"context"
	"go-patron/internal/usecase"
	"go-patron/proto" // Reemplaza por tu módulo real de go.mod
	"log"
	"time"

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
	///////////////////Metadata/////////////////////////////
	// 1. Extraer la metadata entrante del contexto
	// md, exists := metadata.FromIncomingContext(ctx)
	// if !exists {
	// 	return nil, status.Error(codes.Unauthenticated, "No se proporcionó metadata")
	// }
	// // 2. Leer el campo "authorization" (gRPC convierte todas las llaves a MINÚSCULAS automáticamente)
	// authHeader := md.Get("authorization")
	// if len(authHeader) == 0 {
	// 	return nil, status.Error(codes.Unauthenticated, "Token de autorización ausente")
	// }
	// // Extraer el token quitando el prefijo "Bearer "
	// token := strings.TrimPrefix(authHeader[0], "Bearer ")
	// log.Printf("[Servidor] Validando token recibido: %s", token)
	// // Simulación de validación del token
	// if token != "TOKEN" { // Aquí validarías con tu librería JWT usando cfg.JWTKey
	// 	return nil, status.Error(codes.Unauthenticated, "Token inválido o expirado")
	// }
	// //// 3. RESPONDER CON METADATA AL CLIENTE (Headers de salida)
	// headerResponse := metadata.Pairs(
	// 	"x-server-version", "1.0.0",
	// 	"x-request-processed", "true",
	// )
	// // Enviamos el header inmediatamente a través del contexto de gRPC
	// if err := grpc.SendHeader(ctx, headerResponse); err != nil {
	// 	return nil, status.Error(codes.Internal, "No se pudo enviar la metadata de respuesta")
	// }
	////////////////////////////////////////////////

	//Mapear el DTO obtenido del Caso de Uso hacia el UserResponse de Protobuf
	return &proto.UserResponse{
		Id:    int32(userDTO.ID),
		Name:  userDTO.Name,
		Email: userDTO.Email,
	}, nil
}

func (s *Server) GetUsers(req *proto.Empty, stream proto.UserService_GetUsersServer) error {
	log.Println("Cliente conectado solicitando streaming de usuarios...")

	usuarios := s.usecase.FindAll()

	for _, u := range usuarios {
		if err := stream.Context().Err(); err != nil {
			log.Println("El cliente canceló la conexión del stream")
			return status.Error(codes.Canceled, "Conexión cancelada por el cliente")
		}

		res := &proto.UserResponse{
			Id:    int32(u.ID),
			Name:  u.Name,
			Email: u.Email,
		}

		if err := stream.Send(res); err != nil {
			log.Printf("Error al enviar usuario ID %d por el stream: %v", u.ID, err)
			return status.Error(codes.Internal, "Fallo al transmitir datos")
		}

		log.Printf("Usuario %s enviado con éxito", u.Name)

		time.Sleep(1 * time.Second)
	}

	log.Println("Streaming de usuarios finalizado correctamente")
	return nil
}
