package cache

import (
	"backend/config"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func SetHomePageCache(data gin.H) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return config.RedisClient.Set(config.Ctx, HomepageKey(), jsonData, 10*time.Minute).Err()
}

func GetHomePageCache() (gin.H, error) {

	var data gin.H

	cachedData, err := config.RedisClient.Get(config.Ctx, HomepageKey()).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(cachedData), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DeleteHomePageCache() error {

	return config.RedisClient.Del(config.Ctx, HomepageKey()).Err()
}
