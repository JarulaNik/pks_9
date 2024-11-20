package handlers

import (
	"api_note/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetProducts(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		err := db.Select(&products, "SELECT * FROM product")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка получения списка продуктов",
				"details": err.Error(), // Логирование деталей ошибки
			})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProduct(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := strings.TrimSpace(c.Param("id"))
		if idStr == "" {
			log.Println("Параметр id пуст")
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetProduct. ID пользователя отсутствует"})
			return
		}
		log.Printf("GetProduct. Полученный параметр idStr: '%s'", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetProduct. Некорректный ID продукта"})
			return
		}
		var product models.Product
		err = db.Get(&product, "SELECT * FROM Product WHERE product_id = $1", id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "GetProduct. Продукт не найден"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CreateProduct. Некорректные данные"})
			return
		}

		query := `INSERT INTO Product (name, description, price, stock, image_url) 
                  VALUES (:name, :description, :price, :stock, :image_url) RETURNING product_id`

		rows, err := db.NamedQuery(query, &product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "CreateProduct. Ошибка добавления продукта"})
			return
		}
		if rows.Next() {
			rows.Scan(&product.ProductID)
		}
		rows.Close()

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UpdateProduct. Некорректный ID продукта"})
			return
		}

		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UpdateProduct. Некорректные данные"})
			return
		}

		product.ProductID = id
		query := `UPDATE Product SET name = :name, description = :description, price = :price, 
                  stock = :stock, image_url = :image_url WHERE product_id = :product_id`

		_, err = db.NamedExec(query, &product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "UpdateProduct. Ошибка обновления продукта"})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DeleteProduct. Некорректный ID продукта"})
			return
		}

		_, err = db.Exec("DELETE FROM Product WHERE product_id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DeleteProduct. Ошибка удаления продукта"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "DeleteProduct. Продукт успешно удален"})
	}
}
