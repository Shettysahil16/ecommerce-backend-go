package cart

import (
	services "backend/services/cart_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ReduceCartItem(c *gin.Context) {

	productId := c.Param("productId")
	objID, _ := bson.ObjectIDFromHex(productId)

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "please login to continue",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCart, err := services.ReduceCartItem(ctx, userID.(string), objID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cartItem": updatedCart,
	})

}
