package products

import (
	services "backend/services/product_service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func SingleCategoryProduct(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	products, err := services.GetSingleCategoryProduct(ctx)
	//fmt.Print("products in controller", products)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch products",
		})

		return
	}

	c.JSON(200, gin.H{
		"products": products,
	})
}
