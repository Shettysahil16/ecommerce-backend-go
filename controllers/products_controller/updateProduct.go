package products

import (
	services "backend/services/product_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateProduct(c *gin.Context) {

	var req struct {
		ProductID string                 `json:"productID"`
		Products  map[string]interface{} `json:"products"`
	}

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

	product, err := services.UpdateProduct(ctx, userID.(string), req.ProductID, req.Products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated successfully",
		"product": product,
	})

}
