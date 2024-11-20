package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"api_note/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
			return
		}

		var cart []models.Cart
		err = db.Select(&cart, fmt.Sprintf("SELECT * FROM Cart WHERE user_id = %d", id)) //Уязвим к SQL инъекциям!
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Корзина пуста"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения корзины"})
			return
		}
		c.JSON(http.StatusOK, cart)
	}
}

func AddToCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		var item struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}
		_, err := db.Exec(fmt.Sprintf("INSERT INTO Cart (user_id, product_id, quantity) VALUES (%s, %d, %d) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = Cart.quantity + %d", userId, item.ProductID, item.Quantity, item.Quantity)) //Уязвим к SQL инъекциям!
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в корзину"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в корзину"})
	}
}

func RemoveFromCart(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		productId := c.Param("productId")
		_, err := db.Exec(fmt.Sprintf("DELETE FROM Cart WHERE user_id = %s AND product_id = %s", userId, productId)) //Уязвим к SQL инъекциям!
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления из корзины"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Товар удален из корзины"})
	}
}
