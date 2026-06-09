package checkoutcontroller

import (
	cart "backend/services/cart_service"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateCheckoutCartQuantity(c *gin.Context) {

	productId := c.Param("productId")
	productObjId, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, _, ok := utils.GetUser(c)
	if !ok {
		return
	}

	var quantity int
	if err := c.ShouldBindJSON(&quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkoutProduct, err := cart.UpdateCheckoutCartQuantity(ctx, userID, productObjId, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"checkoutProduct": checkoutProduct,
	})
}
