package homepage

import (
	services "backend/services/homepage_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {

	userID := c.GetString("userID")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	data, err := services.GetHomePage(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, data)
}
