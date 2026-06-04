package users

import (
	services "backend/services/user_service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthUser(c *gin.Context) {

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	//start := time.Now()

	user, err := services.GetUserByID(ctx, userID.(string))

	//log.Println("DB query time:", time.Since(start))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})

		return
	}

	//responseStart := time.Now()

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

	//log.Println("Response time:", time.Since(responseStart))

}
