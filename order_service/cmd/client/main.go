package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	userpb "go-patron/proto"
	productpb "product-service/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	_ = godotenv.Load()

	userGRPCURL := getEnv("USER_GRPC_URL", "localhost:50051")
	productGRPCURL := getEnv("PRODUCT_GRPC_URL", "localhost:50052")

	// Conectarse al UserService gRPC.
	userConnection, err := grpc.NewClient(userGRPCURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a UserService: %v", err)
	}
	defer userConnection.Close()

	// Conectarse al ProductService gRPC.
	productConnection, err := grpc.NewClient(productGRPCURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar a ProductService: %v", err)
	}
	defer productConnection.Close()

	userClient := userpb.NewUserServiceClient(userConnection)
	productClient := productpb.NewProductServiceClient(productConnection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	getUser(ctx, userClient, 1)
	//getUsers(ctx, userClient)
	getProduct(ctx, productClient, 1)
	//getProducts(ctx, productClient)
}

func getUser(ctx context.Context, client userpb.UserServiceClient, id int32) {
	user, err := client.GetUser(ctx, &userpb.UserRequest{Id: id})
	if err != nil {
		log.Printf("Error en GetUser: %v", err)
		return
	}

	fmt.Printf("Usuario encontrado: %d - %s (%s)\n", user.GetId(), user.GetName(), user.GetEmail())
}

func getUsers(ctx context.Context, client userpb.UserServiceClient) {
	stream, err := client.GetUsers(ctx, &userpb.Empty{})
	if err != nil {
		log.Printf("Error al iniciar GetUsers: %v", err)
		return
	}

	fmt.Println("Usuarios recibidos por streaming:")
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error al recibir usuario: %v", err)
			return
		}
		fmt.Printf("- %d: %s (%s)\n", user.GetId(), user.GetName(), user.GetEmail())
	}
}

func getProduct(ctx context.Context, client productpb.ProductServiceClient, id int32) {
	product, err := client.GetProduct(ctx, &productpb.ProductRequest{Id: id})
	if err != nil {
		log.Printf("Error en GetProduct: %v", err)
		return
	}

	fmt.Printf("Producto encontrado: %d - %s | precio: %.2f | stock: %d\n", product.GetId(), product.GetName(), product.GetPrice(), product.GetStock())
}

func getProducts(ctx context.Context, client productpb.ProductServiceClient) {
	stream, err := client.GetProducts(ctx, &productpb.Empty{})
	if err != nil {
		log.Printf("Error al iniciar GetProducts: %v", err)
		return
	}

	fmt.Println("Productos recibidos por streaming:")
	for {
		product, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error al recibir producto: %v", err)
			return
		}
		fmt.Printf("- %d: %s | precio: %.2f | stock: %d\n", product.GetId(), product.GetName(), product.GetPrice(), product.GetStock())
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
