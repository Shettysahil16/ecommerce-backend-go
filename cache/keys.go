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
