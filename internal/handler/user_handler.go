package handler

import (
	"go-patron/internal/entity"
	"go-patron/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase *usecase.UserUseCase
}

// Constructor
func NewUserHandler(usecase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

// Método
func (h *UserHandler) Create(c *gin.Context) {

	var user entity.User

	c.BindJSON(&user)

	err := h.usecase.Create(&user)

	if err != nil {

		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(201, user)
}

func (h *UserHandler) FindByID(c *gin.Context) {
	// var userID uint
	// c.BindUri(&userID)
	// 1. Obtener el string desde la URL /user/:id
	idStr := c.Param("id")
	// 2. Convertir el string a un entero
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El ID debe ser un número entero positivo",
		})
		return
	}
	// 3. Castear a uint de forma segura
	userID := uint(idInt)

	user, err := h.usecase.FindByID(userID)

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, user)
}

//El Handler solo recibe datos y responde. Toda la lógica está en el Use Case.
