package main

import (
	"api_note/db"
	"api_note/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаемся к базе данных
	connectDB, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	router := gin.Default()

	// Роуты для продуктов
	router.GET("/products", handlers.GetProducts(connectDB))
	router.GET("/products/:id", handlers.GetProduct(connectDB))
	router.POST("/products", handlers.CreateProduct(connectDB))
	router.PUT("/products/:id", handlers.UpdateProduct(connectDB))
	router.DELETE("/products/:id", handlers.DeleteProduct(connectDB))

	// Роуты для корзины
	router.GET("/carts/:id", handlers.GetCart(connectDB))
	router.POST("/carts/:userId", handlers.AddToCart(connectDB))
	router.DELETE("/carts/:userId/:productId", handlers.RemoveFromCart(connectDB))

	// Роуты для избранного
	router.GET("/favorites/:id", handlers.GetFavorites(connectDB))
	router.POST("/favorites/:userId", handlers.AddToFavorites(connectDB))
	router.DELETE("/favorites/:userId/:productId", handlers.RemoveFromFavorites(connectDB))

	// Роуты для заказов
	router.GET("/orders/:id", handlers.GetOrders(connectDB))
	router.POST("/orders/:id", handlers.CreateOrder(connectDB))

	// Запуск сервера
	err = router.Run(":8080")
	if err != nil {
		return
	}
}
