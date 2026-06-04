package users

import (
	"backend/models"
	services "backend/services/user_service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var user models.LoginRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, err := services.LoginUser(ctx, user)
	if err != nil {
		c.JSON(401, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", false, true)

	c.JSON(200, gin.H{
		"message":     "login successful",
		"accessToken": accessToken,
	})
}
