package cache

import (
	"backend/config"
	"backend/models"
	"context"
)

func GetCheckout(ctx context.Context, userID string) (map[string]string, error) {

	data, err := config.RedisClient.HGetAll(ctx, UserCheckoutKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SetCheckoutItems(ctx context.Context, userID string, items []models.CheckoutItem) error {

	pipe := config.RedisClient.TxPipeline()
	pipe.Del(ctx, UserCheckoutKey(userID))

	if len(items) > 0 {
		data := make(map[string]any, len(items))

		for _, item := range items {
			data[item.ProductID] = item.Quantity
		}

		pipe.HSet(ctx, UserCheckoutKey(userID), data)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func UpdateCheckoutItem(ctx context.Context, userID string, productID string, quantity int64) error {
	return config.RedisClient.HSet(
		ctx,
		UserCheckoutKey(userID),
		productID,
		quantity,
	).Err()
}

func RemoveCheckoutItem(ctx context.Context, userID string, productID string) error {
	return config.RedisClient.HDel(
		ctx,
		UserCheckoutKey(userID),
		productID,
	).Err()
}
