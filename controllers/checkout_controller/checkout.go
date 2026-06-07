package checkoutcontroller

import (
	"backend/models"
	checkoutService "backend/services/checkout_service"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
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

	result, err := checkoutService.CheckoutService(ctx, req, userID, userObjId)
	if err != nil {
		checkoutService.HandleCheckoutError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
