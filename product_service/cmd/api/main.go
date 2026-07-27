package main

import (
	"fmt"
	"log"
	"net"
	"product-service/internal/config"
	"product-service/internal/entity"
	grpc_server "product-service/internal/grpc"
	"product-service/internal/handler"
	"product-service/internal/repository"
	"product-service/internal/usecase"
	productpb "product-service/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.BDName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	// Crea la tabla si todavía no existe. Es útil para las primeras clases.
	if err := db.AutoMigrate(&entity.Product{}); err != nil {
		log.Fatalf("No se pudo crear la tabla de productos: %v", err)
	}

	productRepo := repository.NewProductRepository(db)
	productUseCase := usecase.NewProductUseCase(productRepo)
	productHandler := handler.NewProductHandler(productUseCase)

	// Inicializar el servidor gRPC usando el mismo caso de uso que la API REST.
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("No se pudo abrir el puerto gRPC: %v", err)
	}
	grpcServer := grpc.NewServer()
	productGRPCServer := grpc_server.NewProductGRPCServer(productUseCase)
	productpb.RegisterProductServiceServer(grpcServer, productGRPCServer)

	go func() {
		log.Printf("ProductService gRPC ejecutándose en :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Error en el servidor gRPC: %v", err)
		}
	}()

	r := gin.Default()
	r.GET("/products", productHandler.FindAll)
	r.GET("/product/:id", productHandler.FindByID)
	r.POST("/product", productHandler.Create)
	r.PUT("/product/:id", productHandler.Update)
	r.DELETE("/product/:id", productHandler.Delete)

	log.Printf("product_service ejecutándose en :%s", cfg.APIPort)
	if err := r.Run(":" + cfg.APIPort); err != nil {
		log.Fatal(err)
	}
}

//protoc --go_out=. --go-grpc_out=. user.proto
