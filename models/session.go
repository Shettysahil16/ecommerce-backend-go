package models

type Session struct {
	SessionID   string `json:"sessionId"`
	UserID      string `json:"userId"`
	RefreshHash string `json:"refreshHash"`
}
