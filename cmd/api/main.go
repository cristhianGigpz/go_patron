package main

import (
	"fmt"
	"go-patron/internal/config"
	grpc_server "go-patron/internal/grpc"
	"go-patron/internal/handler"
	"go-patron/internal/repository"
	"go-patron/internal/usecase"
	"go-patron/proto"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//const dsn = "host=localhost user=postgres password=gigpz dbname=bd_tests port=5434 sslmode=disable"

func main() {
	fmt.Println("Hello, World! Clean Architecture in Go")
	cfg := config.LoadConfig()

	// if cfg.JWTKey == "" {
	// 	log.Fatal("Error crítico: La variable JWT_KEY es obligatoria y no está configurada")
	// }

	// 1. Inicializar la conexión a la base de datos
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.BDName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	r := gin.Default()
	// 2. Inyección de Dependencias (Tu código)
	// Pasamos 'db' al repositorio
	repo := repository.NewUserRepository(db)
	//repo := mock.MockRepository{}

	// Pasamos 'repo' al caso de uso
	//userUseCase := usecase.NewUserUseCase(&repo)
	userUseCase := usecase.NewUserUseCase(repo)
	// Inyectamos el caso de uso en nuestro servidor gRPC///
	grpcUserService := grpc_server.NewUserGRPCServer(userUseCase)
	// 4. Iniciar la escucha en el puerto de red
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fallo al escuchar en puerto 50051: %v", err)
	}
	// 5. Configurar y encender el servidor gRPC nativo

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_server.AuthUnaryInterceptor),
	)
	// Registramos nuestro servicio inyectado
	proto.RegisterUserServiceServer(server, grpcUserService)

	log.Println("Microservicio gRPC corriendo en el puerto :50051 con acceso real a Base de Datos")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Error al levantar el servidor gRPC: %v", err)
	}

	// Pasamos 'userUseCase' al manejador
	handler := handler.NewUserHandler(userUseCase)
	// 3. A partir de aquí, utilizas 'userUseCase'
	// Por ejemplo, pasándolo a tus controladores HTTP (Gin, Fiber, etc.)
	// o ejecutando una prueba rápida:

	r.GET("/user/:id", func(c *gin.Context) {
		handler.FindByID(c)
	})

	r.POST("/user", func(c *gin.Context) {
		handler.Create(c)
	})

	r.Run(":8080")

	//log.Println("Aplicación inicializada correctamente con Inyección de Dependencias", userUseCase)
	// 				Cliente
	//              │
	//              ▼
	//       Gin (Handlers)
	//              │
	//              ▼
	//         Use Cases
	//              │
	//              ▼
	//   Repository Interfaces
	//              │
	//              ▼
	// GORM / Redis / APIs externas
	//              │
	//              ▼
	//    PostgreSQL / Redis

	/////////////////gRPC/////////////////////////
	// 	Cliente
	//    │
	//    ▼
	// Stub generado por Protobuf
	//    │
	//    ▼
	// HTTP/2 + Protobuf
	//    │
	//    ▼
	// Servidor gRPC
	//    │
	//    ▼
	// Servicio
	//    │
	//    ▼
	// Base de Datos
}
