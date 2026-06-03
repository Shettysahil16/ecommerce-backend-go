package models

type CartItemResponse struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}
