package main

import (
	"fmt"
	"go-patron/internal/config"
	"go-patron/internal/handler"
	"go-patron/internal/repository"
	"go-patron/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
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
}
