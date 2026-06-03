package cart

import (
	services "backend/services/cart_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCartItems(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "please login to continue",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cartItems, err := services.GetCartItems(ctx, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(cartItems.Items) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No products found in cart",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart": cartItems,
	})

}
