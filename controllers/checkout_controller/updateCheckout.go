package checkoutcontroller

import (
	checkoutservice "backend/services/checkout_service"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UpdateCheckoutController(c *gin.Context) {

	productId := c.Param("productId")
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, _, ok := utils.GetUser(c)
	if !ok {
		return
	}

	var qty int64
	if err := c.ShouldBindJSON(&qty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := checkoutservice.UpdateCheckoutItemService(ctx, userID, productId, qty)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no product found",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "prouduct quantity updated",
	})
}
