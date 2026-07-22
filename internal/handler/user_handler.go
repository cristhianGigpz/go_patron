package handler

import (
	"go-patron/internal/dto"
	"go-patron/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase usecase.UserUseCaseInterface
}

// Constructor
func NewUserHandler(usecase usecase.UserUseCaseInterface) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

// Método
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	//var user entity.User
	//c.BindJSON(&user)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Enviamos el DTO al caso de uso y recibimos el DTO de respuesta seguro
	res, err := h.usecase.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna un 201 con los datos limpios (sin password)
	c.JSON(http.StatusCreated, res)
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

	userDTO, err := h.usecase.FindByID(userID)
	if err != nil {
		// Aquí puedes mapear errores específicos de GORM (como registro no encontrado) a un 404
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Retorna el DTO sin campos ocultos o contraseñas
	c.JSON(http.StatusOK, userDTO)
}

func (h *UserHandler) FindAll(c *gin.Context) {
	users := h.usecase.FindAll()
	c.JSON(http.StatusOK, users)
}

//El Handler solo recibe datos y responde. Toda la lógica está en el Use Case.
