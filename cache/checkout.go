package cache

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"time"
)

func SetCheckoutProducts(userID string, checkoutProduct *models.CheckoutResponse) error {

	jsonData, err := json.Marshal(checkoutProduct)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(config.Ctx, UserCheckoutKey(userID), jsonData, 10*time.Minute).Err()
}
