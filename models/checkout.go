package models

type CheckoutRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	AddressID string `json:"addressId" binding:"required"`
}

type CheckoutResponse struct {
	Items      []CheckoutItemResponse  `json:"items"`
	TotalPrice float64                 `json:"totalPrice"`
	Address    *DefaultAddressResponse `json:"address"`
}

type CheckoutItemResponse struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type CheckoutItem struct {
	ProductID string `json:"productId"`
	Quantity  int64  `json:"quantity"`
}

type CheckoutCache struct {
	Items []CheckoutItem `json:"items"`
}
