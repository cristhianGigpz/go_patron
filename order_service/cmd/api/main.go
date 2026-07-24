package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func main() {
	_ = godotenv.Load()
	userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	productServiceURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8081")

	users, err := getUsers(userServiceURL)
	if err != nil {
		log.Printf("No se pudieron obtener los usuarios: %v", err)
	} else {
		fmt.Println("Usuarios:")
		for _, user := range users {
			fmt.Printf("- %d: %s (%s)\n", user.ID, user.Name, user.Email)
		}
	}

	products, err := getProducts(productServiceURL)
	if err != nil {
		log.Printf("No se pudieron obtener los productos: %v", err)
	} else {
		fmt.Println("Productos:")
		for _, product := range products {
			fmt.Printf("- %d: %s | precio: %.2f | stock: %d\n", product.ID, product.Name, product.Price, product.Stock)
		}
	}
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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
