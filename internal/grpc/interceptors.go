package grpc_server

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthUnaryInterceptor valida el token JWT/Bearer de forma global
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("[Interceptor] Interceptando petición para el método: %s", info.FullMethod)

	// 1. Extraer la metadata del contexto entrante
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		return nil, status.Error(codes.Unauthenticated, "Metadata requerida ausente")
	}

	// 2. Leer el campo 'authorization'
	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Token de autorización ausente")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	// 3. Validar el Token (Simulación)
	if token != "TOKEN" {
		log.Printf("[Interceptor] Acceso denegado para el método %s", info.FullMethod)
		return nil, status.Error(codes.Unauthenticated, "Token inválido o expirado")
	}

	log.Println("[Interceptor] Token válido. Continuando con el servicio...")

	// 4. Continuar con la ejecución normal del método del servidor (ej. GetUser)
	return handler(ctx, req)
}
