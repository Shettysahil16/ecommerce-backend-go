package cache

import (
	"backend/config"
	"backend/models"
)

func AddToCartCache(userID string, productID string, qty int) (*models.RedisCart, error) {

	newQty, err := config.RedisClient.HIncrBy(config.Ctx, CartKey(userID), productID, int64(qty)).Result()
	if err != nil {
		return nil, err
	}

	return &models.RedisCart{
		ProductID: productID,
		Quantity:  newQty,
	}, nil

}

func GetCartCache(userID string) (map[string]string, error) {

	data, err := config.RedisClient.HGetAll(config.Ctx, CartKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return data, nil

}

func ReduceCartCache(userID string, productID string, qty int) (*models.RedisCart, error) {
	newQty, err := config.RedisClient.HIncrBy(config.Ctx, CartKey(userID), productID, int64(-qty)).Result()
	if err != nil {
		return nil, err
	}

	if newQty <= 0 {
		err := config.RedisClient.HDel(config.Ctx, CartKey(userID), productID).Err()
		if err != nil {
			return nil, err
		}

		newQty = 0

	}

	return &models.RedisCart{
		ProductID: productID,
		Quantity:  newQty,
	}, nil
}

func DeleteCartItemCache(userID string, productID string) error {

	return config.RedisClient.HDel(config.Ctx, CartKey(userID), productID).Err()
}

func UpdateCartQuantityCache(userID string, productID string, qty int) (*models.RedisCart, error) {

	newQty, err := config.RedisClient.HSet(config.Ctx, CartKey(userID), productID, int64(qty)).Result()
	if err != nil {
		return nil, err
	}

	return &models.RedisCart{
		ProductID: productID,
		Quantity:  newQty,
	}, nil
}
