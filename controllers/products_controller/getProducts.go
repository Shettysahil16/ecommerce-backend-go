package products

import (
	services "backend/services/product_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	products, err := services.GetProducts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
