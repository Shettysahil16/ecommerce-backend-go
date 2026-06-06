package controllers

import (
	"backend/cache"
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {

	// READING REFRESH TOKEN FROM USER COOKIE
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token missing",
		})
		return
	}

	// VALIDATING THE USER REFRESH TOKEN
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		},
		)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// FINDING THE SESSION INSIDE THE REDIS
	session, err := cache.GetSessionCache(ctx, claims.SessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "session not found",
		})
		return
	}

	// HASHING THE REFRESH TOKEN COMING FROM USER
	incomingHash := utils.HashToken(refreshToken)

	// CHECKING WHETHER THE BOTH HASHED REFRESH TOKEN FROM CURRENT REQUEST AND STORED REDIS REFRESH TOKEN ARE SAME OR NOT
	if incomingHash != session.RefreshHash {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(session.UserID, session.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate refresh token",
		})
		return
	}

	newRefreshHash := utils.HashToken(newRefreshToken)

	// CHANGING SESSION KEY REFRESHHASH WITH NEW REFRESHHASH MADE
	// THIS IS SHORT FORM BECAUSE THE ONLY VALUE OF REFRESHHASH IS CHANGING REST OTHER ARE SAME
	session.RefreshHash = newRefreshHash

	err = cache.SetSessionCache(ctx, claims.SessionID, *session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to store refresh token",
		})
		return
	}

	c.SetCookie("refresh_token", newRefreshToken, 7*24*60*60, "/", "", false, true)

	// AFTER THE SUCCESS GENERATING THE NEW ACCESS TOKEN
	accessToken, err := utils.GenerateAcccessToken(claims.UserID, claims.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	},
	)

}
