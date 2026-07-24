package main

import (
	"fmt"
	"log"
	"product-service/internal/config"
	"product-service/internal/entity"

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

	if err := db.AutoMigrate(&entity.Product{}); err != nil {
		log.Fatalf("No se pudo crear la tabla products: %v", err)
	}

	var count int64
	db.Model(&entity.Product{}).Count(&count)
	if count > 0 {
		log.Printf("La tabla products ya tiene %d registros; no se insertaron duplicados", count)
		return
	}

	products := []entity.Product{
		{Name: "Laptop Lenovo IdeaPad", Description: "Laptop para trabajo y estudio", Price: 2499.90, Stock: 12},
		{Name: "Mouse Logitech M185", Description: "Mouse inalámbrico", Price: 59.90, Stock: 35},
		{Name: "Teclado Mecánico Redragon", Description: "Teclado mecánico RGB", Price: 189.90, Stock: 18},
		{Name: "Monitor LG 24 pulgadas", Description: "Monitor Full HD", Price: 699.90, Stock: 10},
		{Name: "Audífonos Sony WH-CH520", Description: "Audífonos Bluetooth", Price: 179.90, Stock: 22},
		{Name: "Webcam Logitech C920", Description: "Cámara Full HD para videollamadas", Price: 329.90, Stock: 8},
		{Name: "Disco SSD Kingston 480GB", Description: "Unidad de almacenamiento SSD", Price: 159.90, Stock: 25},
		{Name: "Memoria USB Kingston 64GB", Description: "Memoria USB 3.0", Price: 29.90, Stock: 50},
		{Name: "Router TP-Link Archer C6", Description: "Router Wi-Fi de doble banda", Price: 219.90, Stock: 14},
		{Name: "Cable HDMI 2 metros", Description: "Cable HDMI de alta velocidad", Price: 24.90, Stock: 40},
		{Name: "Tablet Samsung Galaxy Tab A9", Description: "Tablet de 8.7 pulgadas", Price: 749.90, Stock: 7},
		{Name: "Smartphone Motorola G54", Description: "Teléfono móvil 5G", Price: 899.90, Stock: 9},
		{Name: "Cargador USB-C 25W", Description: "Cargador de carga rápida", Price: 69.90, Stock: 30},
		{Name: "Power Bank Xiaomi 10000mAh", Description: "Batería portátil", Price: 99.90, Stock: 16},
		{Name: "Parlante JBL Go 3", Description: "Parlante Bluetooth portátil", Price: 159.90, Stock: 11},
		{Name: "Silla Ergonómica", Description: "Silla para oficina", Price: 549.90, Stock: 5},
		{Name: "Escritorio 120cm", Description: "Escritorio de melamina", Price: 399.90, Stock: 6},
		{Name: "Mochila para Laptop", Description: "Mochila impermeable", Price: 119.90, Stock: 20},
		{Name: "Impresora Multifuncional Epson", Description: "Impresora con conexión Wi-Fi", Price: 629.90, Stock: 4},
		{Name: "Regulador de Voltaje", Description: "Regulador de voltaje de 1000VA", Price: 129.90, Stock: 13},
	}

	if err := db.Create(&products).Error; err != nil {
		log.Fatalf("No se pudieron insertar los productos: %v", err)
	}

	log.Printf("Se insertaron %d productos en la tabla products", len(products))
}
