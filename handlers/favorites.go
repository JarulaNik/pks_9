package handlers

import (
	"api_note/models"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

func errorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func GetFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			errorResponse(c, http.StatusBadRequest, "Неверный формат ID пользователя: "+userIDStr)
			return
		}

		var favorites []models.Favorite
		err = db.Select(&favorites, "SELECT * FROM Favorites WHERE user_id = ?", userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				errorResponse(c, http.StatusNotFound, "Список избранного пуст")
				return
			}
			errorResponse(c, http.StatusInternalServerError, "Ошибка получения списка избранного: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, favorites)
	}
}

func AddToFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			UserID    int `uri:"userId" binding:"required"`
			ProductID int `json:"product_id" binding:"required,min=1"`
		}
		if err := c.ShouldBindUri(&request); err != nil {
			errorResponse(c, http.StatusBadRequest, "Неверный формат ID пользователя")
			return
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			errorResponse(c, http.StatusBadRequest, "Неверный формат данных: "+err.Error())
			return
		}

		res, err := db.Exec("INSERT INTO Favorites (user_id, product_id) VALUES (?, ?) ON CONFLICT DO NOTHING", request.UserID, request.ProductID)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Ошибка добавления в избранное: "+err.Error())
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Ошибка получения количества измененных строк: "+err.Error())
			return
		}

		if rowsAffected > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в избранное"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Товар уже есть в избранном"})
		}
	}
}

func RemoveFromFavorites(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			UserID    int `uri:"userId" binding:"required,min=1"`
			ProductID int `uri:"productId" binding:"required,min=1"`
		}
		if err := c.ShouldBindUri(&request); err != nil {
			errorResponse(c, http.StatusBadRequest, "Неверный формат ID пользователя или товара: "+err.Error())
			return
		}

		res, err := db.Exec("DELETE FROM Favorites WHERE user_id = ? AND product_id = ?", request.UserID, request.ProductID)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Ошибка удаления из избранного: "+err.Error())
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Ошибка подсчета удаленных строк: "+err.Error())
			return
		}

		if rowsAffected > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Товар удален из избранного"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден в избранном"})
		}
	}
}
