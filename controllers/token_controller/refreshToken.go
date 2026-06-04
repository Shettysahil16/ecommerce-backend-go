package controllers

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {

	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		},
		)
		return
	}

	accessToken, err := utils.GenerateAcccessToken(claims.UserID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	},
	)

}
