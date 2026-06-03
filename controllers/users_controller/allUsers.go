package users

import (
	services "backend/services/user_service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func AllUsers(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, err := services.GetAllUsers(ctx)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})

		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
