package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-patron/proto" // Reemplaza por tu módulo real de go.mod

	"github.com/gin-gonic/gin"
)

type UserWebHandlerRpc struct {
	// Inyectamos el cliente gRPC generado por protobuf
	grpcClient proto.UserServiceClient
}

// Constructor del Handler
func NewUserWebHandlerRpc(client proto.UserServiceClient) *UserWebHandlerRpc {
	return &UserWebHandlerRpc{
		grpcClient: client,
	}
}

func (h *UserWebHandlerRpc) FindByID(c *gin.Context) {
	// 1. Validar el ID que viene por la URL (HTTP GET /users/:id)
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero positivo"})
		return
	}

	// 2. Crear un contexto con Timeout para la llamada gRPC interna
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 3. Preparar el Request para el microservicio gRPC
	grpcReq := &proto.UserRequest{
		Id: int32(idInt), // Convertimos al tipo de dato que espera tu proto
	}

	// 4. Realizar la llamada interna por gRPC al microservicio
	grpcRes, err := h.grpcClient.GetUser(ctx, grpcReq)
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
}
