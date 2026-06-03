package products

import (
	"backend/models"
	services "backend/services/product_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadProduct(c *gin.Context) {

	var req models.Product

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	product, err := services.UploadProduct(ctx, userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"product": product,
	})

}
