package cart

import (
	services "backend/services/cart_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddToCart(c *gin.Context) {

	productId := c.Param("productId")
	objID, _ := bson.ObjectIDFromHex(productId)
	//fmt.Println("productId", objID)

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "please login to continue",
		})
		return
	}
	//fmt.Println("user id", userID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	updatedCart, err := services.AddToCart(ctx, userID.(string), objID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//fmt.Println("cartItem", cartItem)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "product added to cart",
		"updatedCart": updatedCart,
	})
}
