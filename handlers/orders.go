package handlers

import (
	"api_note/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetOrders(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID из параметров маршрута
		idStr := c.Param("id")
		log.Println("GetOrders. Полученный параметр idStr:", idStr)

		// Убираем лишние пробелы
		idStr = strings.TrimSpace(idStr)

		// Преобразуем ID в целое число
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("Ошибка преобразования idStr в int:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
			return
		}

		// Выполняем запрос к базе данных
		var orders []models.Order
		err = db.Select(&orders, "SELECT * FROM orders WHERE user_id = $1", id)

		if err != nil {
			log.Println("Ошибка запроса к базе данных:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения заказов"})
			return
		}

		// Отправляем ответ
		c.JSON(http.StatusOK, orders)
	}
}

func CreateOrder(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
			return
		}

		query := `
            INSERT INTO orders (user_id, total, status, created_at)
            VALUES (:user_id, :total, :status, :created_at)
        `

		rows, err := db.NamedQuery(query, &order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления заказа"})
			return
		}
		if rows.Next() {
			err := rows.Scan(&order.OrderID)
			if err != nil {
				return
			}
		}
		err = rows.Close()
		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, order)
	}
}
