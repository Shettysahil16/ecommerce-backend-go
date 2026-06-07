package checkoutservice

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ErrCartEmpty = errors.New("cart is empty")

func HandleCheckoutError(c *gin.Context, err error) {
	if errors.Is(err, ErrCartEmpty) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cart is empty",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
