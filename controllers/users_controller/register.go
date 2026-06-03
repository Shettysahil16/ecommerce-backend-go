package users

import (
	"backend/models"
	services "backend/services/user_service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = services.RegisterUser(ctx, user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "user created successfully",
	})
}
