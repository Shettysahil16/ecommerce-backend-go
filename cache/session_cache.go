package cache

import (
	"backend/config"
	"backend/models"
	"context"
	"encoding/json"
	"time"
)

func SetSessionCache(ctx context.Context, sessionID string, session models.Session) error {

	jsonData, err := json.Marshal(session)
	if err != nil {
		return err
	}

	pipe := config.RedisClient.TxPipeline()

	// CURRENT SESSION
	pipe.Set(ctx, SessionKey(sessionID), jsonData, 7*24*time.Hour)

	// ALL SESSIONS OF SAME USER
	pipe.SAdd(ctx, UserSessionsKey(session.UserID), sessionID)

	_, err = pipe.Exec(ctx)

	return err
}

func GetUserSessionIDS(ctx context.Context, userID string) ([]string, error) {

	return config.RedisClient.SMembers(ctx, UserSessionsKey(userID)).Result()

}

func GetSessionCache(ctx context.Context, sessionID string) (*models.Session, error) {

	data, err := config.RedisClient.Get(ctx, SessionKey(sessionID)).Result()
	if err != nil {
		return nil, err
	}

	var session models.Session

	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetUserSessionsCache(ctx context.Context, userID string) ([]models.Session, error) {

	sessionIDS, err := GetUserSessionIDS(ctx, userID)
	if err != nil {
		return nil, err
	}

	var sessions []models.Session

	for _, sessionID := range sessionIDS {
		session, err := GetSessionCache(ctx, sessionID)
		if err != nil {
			continue
		}

		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func DeleteSessionCache(ctx context.Context, sessionID string) error {

	session, err := GetSessionCache(ctx, sessionID)
	if err != nil {
		return err
	}

	pipe := config.RedisClient.TxPipeline()

	pipe.Del(ctx, SessionKey(sessionID))

	pipe.SRem(ctx, UserSessionsKey(session.UserID), sessionID)

	_, err = pipe.Exec(ctx)

	return err
}

func DeleteAllUsersCache(ctx context.Context, userID string) error {

	sessionIDS, err := GetUserSessionIDS(ctx, userID)
	if err != nil {
		return err
	}

	pipe := config.RedisClient.TxPipeline()

	for _, sessionID := range sessionIDS {
		pipe.Del(ctx, SessionKey(sessionID))
	}

	pipe.Del(ctx, UserSessionsKey(userID))

	_, err = pipe.Exec(ctx)

	return err
}

func IsSessionActive(sessionID string) bool {

	exists, err := config.RedisClient.Exists(config.Ctx, SessionKey(sessionID)).Result()
	if err != nil {
		return false
	}

	return exists == 1
}
