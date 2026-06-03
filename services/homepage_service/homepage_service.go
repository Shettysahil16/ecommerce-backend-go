package services

import (
	"backend/cache"
	product "backend/services/product_service"
	user "backend/services/user_service"
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetHomePage(ctx context.Context, userID string) (gin.H, error) {

	var response gin.H

	homepageCache, err := cache.GetHomePageCache()
	if err == nil {

		fmt.Println("HOMEPAGE CACHE HIT")

		response = homepageCache

	} else {

		fmt.Println("HOMEPAGE CACHE MISS")

		var wg sync.WaitGroup

		var categoryPreview []bson.M
		var airpodes []bson.M

		var previewErr error
		var airpodesErr error

		wg.Add(1)

		go func() {
			defer wg.Done()

			categoryPreview, previewErr = product.GetSingleCategoryProduct(ctx)
		}()

		wg.Add(1)

		go func() {
			defer wg.Done()

			airpodes, airpodesErr = product.GetCategoryProducts(ctx, "airpodes")
		}()

		wg.Wait()

		if previewErr != nil {
			return nil, previewErr
		}

		if airpodesErr != nil {
			return nil, airpodesErr
		}

		response = gin.H{
			"categoryPreview": categoryPreview,
			"airpodes":        airpodes,
		}

		cache.SetHomePageCache(response)
	}

	if userID != "" {

		cachedUser, err := cache.GetUserCache(userID)

		if err == nil {

			fmt.Println("USER CACHE HIT")

			response["user"] = cachedUser
		} else {
			fmt.Println("USER CACHE MISS")

			user, err := user.GetUserByID(ctx, userID)

			if err != nil {
				return nil, err
			}

			cache.SetUserCache(userID, user)

			response["user"] = user
		}
	}

	return response, nil
}
