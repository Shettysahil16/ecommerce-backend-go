package cache

func HomepageKey() string {
	return "homepage"
}
func UserKey(userID string) string {
	return "user:" + userID
}

func CartKey(userID string) string {
	return "cart:" + userID
}

func SessionKey(sessionID string) string {
	return "session:" + sessionID
}

func UserSessionsKey(userID string) string {
	return "user_sessions:" + userID
}
