package models

type RedisCart struct {
	ProductID string `json:"productId"`
	Quantity  int64  `json:"quantity"`
}
