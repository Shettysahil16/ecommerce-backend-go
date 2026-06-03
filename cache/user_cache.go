package cache

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"time"
)

func SetUserCache(userID string, user *models.User) error {

	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(config.Ctx, UserKey(userID), jsonData, 10*time.Minute).Err()
}

func GetUserCache(userID string) (models.User, error) {

	var user models.User

	cachedData, err := config.RedisClient.Get(config.Ctx, UserKey(userID)).Result()
	if err != nil {
		return user, err
	}

	err = json.Unmarshal([]byte(cachedData), &user)

	return user, err
}

func DeleteUserCache(userID string) error {

	return config.RedisClient.Del(config.Ctx, UserKey(userID)).Err()
}
