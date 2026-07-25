package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Product struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type CreateOrderRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderResponse struct {
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
}

func main() {
	_ = godotenv.Load()
	userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	productServiceURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")

	r := gin.Default()

	// Estas rutas permiten consultar los servicios desde el cliente.
	r.GET("/users", func(c *gin.Context) {
		users, err := getUsers(userServiceURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.GET("/products", func(c *gin.Context) {
		products, err := getProducts(productServiceURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// Esta es la operación principal del API Gateway.
	r.POST("/orders", func(c *gin.Context) {
		createOrder(c, userServiceURL, productServiceURL)
	})

	port := getEnv("API_PORT", "8082")
	log.Printf("order_service API Gateway ejecutándose en :%s", port)
	log.Fatal(r.Run(":" + port))
}

func createOrder(c *gin.Context, userServiceURL, productServiceURL string) {
	var request CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el cuerpo debe ser un JSON válido"})
		return
	}

	if request.UserID == 0 || request.ProductID == 0 || request.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, product_id y quantity son obligatorios; quantity debe ser mayor que cero"})
		return
	}

	user, err := getUser(userServiceURL, request.UserID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no se pudo validar el usuario: " + err.Error()})
		return
	}

	product, err := getProduct(productServiceURL, request.ProductID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "no se pudo validar el producto: " + err.Error()})
		return
	}

	if product.Stock < request.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no hay stock suficiente para el producto"})
		return
	}

	order := OrderResponse{
		User:     user,
		Product:  product,
		Quantity: request.Quantity,
		Total:    product.Price * float64(request.Quantity),
	}
	c.JSON(http.StatusCreated, order)
}

func getUser(userServiceURL string, id uint) (User, error) {
	var user User
	err := getJSON(userServiceURL+"/user/"+strconv.FormatUint(uint64(id), 10), &user)
	return user, err
}

func getProduct(productServiceURL string, id uint) (Product, error) {
	var product Product
	err := getJSON(productServiceURL+"/product/"+strconv.FormatUint(uint64(id), 10), &product)
	return product, err
}

func getUsers(userServiceURL string) ([]User, error) {
	response, err := http.Get(userServiceURL + "/users")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("%s respondió %s: %s", userServiceURL, response.Status, string(body))
	}

	var users []User
	if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func getProducts(productServiceURL string) ([]Product, error) {
	response, err := http.Get(productServiceURL + "/products")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("%s respondió %s: %s", productServiceURL, response.Status, string(body))
	}

	var products []Product
	if err := json.NewDecoder(response.Body).Decode(&products); err != nil {
		return nil, err
	}
	return products, nil
}

func getJSON(url string, result interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("%s respondió %s: %s", url, response.Status, string(body))
	}

	return json.NewDecoder(response.Body).Decode(result)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
