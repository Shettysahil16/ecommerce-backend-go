package models

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}
