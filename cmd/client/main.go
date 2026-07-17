package main

import (
	"context"
	"log"
	"time"

	"go-patron/proto" // Reemplaza "go-patron" por el nombre de tu módulo en go.mod

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Establecer la conexión con el servidor gRPC
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor gRPC: %v", err)
	}
	defer conn.Close() // Nos aseguramos de cerrar la conexión al finalizar

	// 2. Crear el cliente utilizando el código generado por proto
	client := proto.NewUserServiceClient(conn)

	// 3. Crear un contexto con tiempo límite (Timeout) para la petición (Buena práctica)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 4. Construir la petición (Request)
	req := &proto.UserRequest{
		Id: 10, // Enviamos el ID que queremos buscar
	}

	// 5. Ejecutar la llamada RPC al servidor //Consumir servicio:
	log.Println("Enviando petición gRPC al servidor...")
	res, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("Error al llamar a GetUser: %v", err)
	}

	// 6. Procesar y mostrar la respuesta del servidor //Resultado:
	log.Printf("Respuesta recibida con éxito del Servidor:")
	log.Printf("ID: %d", res.GetId())
	log.Printf("Nombre: %s", res.GetName())
	log.Printf("Email: %s", res.GetEmail())
}
