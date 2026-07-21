package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-patron/proto" // Reemplaza "go-patron" por el nombre de tu módulo en go.mod

	"github.com/gin-gonic/gin"
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
	//////////////////////////////////////////////////////////

	// // 3. Crear un contexto con tiempo límite (Timeout) para la petición (Buena práctica)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()

	// // 4. Construir la petición (Request)
	// req := &proto.UserRequest{
	// 	Id: 4, // Enviamos el ID que queremos buscar
	// }

	// // 5. Ejecutar la llamada RPC al servidor //Consumir servicio:
	// log.Println("Enviando petición gRPC al servidor...")
	// res, err := client.GetUser(ctx, req)
	// if err != nil {
	// 	log.Fatalf("Error al llamar a GetUser: %v", err)
	// }

	// // 6. Procesar y mostrar la respuesta del servidor //Resultado:
	// log.Printf("Respuesta recibida con éxito del Servidor:")
	// log.Printf("ID: %d", res.GetId())
	// log.Printf("Nombre: %s", res.GetName())
	// log.Printf("Email: %s", res.GetEmail())

	//////////////////////////////////////////////////////////
	// userHandler := handler.NewUserWebHandlerRpc(client)

	// // 4. Configurar el servidor HTTP con Gin
	r := gin.Default()

	r.GET("/user/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		idInt, err := strconv.Atoi(idStr)
		if err != nil || idInt <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero positivo"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		grpcReq := &proto.UserRequest{
			Id: int32(idInt), // Convertimos al tipo de dato que espera tu proto
		}

		// 4. Realizar la llamada interna por gRPC al microservicio
		grpcRes, err := client.GetUser(ctx, grpcReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No se pudo obtener el usuario desde el servicio interno",
				"details": err.Error(),
			})
			return
		}

		// 5. Responder al cliente web mapeando los datos gRPC a JSON
		c.JSON(http.StatusOK, gin.H{
			"id":    grpcRes.GetId(),
			"name":  grpcRes.GetName(),
			"email": grpcRes.GetEmail(),
		})
	})

	log.Println("API HTTP Gateway corriendo en el puerto :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("No se pudo encender el servidor Gin: %v", err)
	}

	// // Definimos la ruta HTTP que usará el usuario final
	// r.GET("/user/:id", userHandler.FindByID)

	// log.Println("API HTTP Gateway corriendo en el puerto :8080")
	// if err := r.Run(":8080"); err != nil {
	// 	log.Fatalf("No se pudo encender el servidor Gin: %v", err)
	// }
}
