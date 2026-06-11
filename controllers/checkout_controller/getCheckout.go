package checkoutcontroller

import (
	"backend/models"
	checkoutservice "backend/services/checkout_service"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCheckoutController(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, userObjId, ok := utils.GetUser(c)
	if !ok {
		return
	}

	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkoutResponse, err := checkoutservice.PrepareCheckout(ctx, userID, userObjId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"checkoutResponse": checkoutResponse,
	})
}
