package models

type CheckoutRequest struct {
	ProductID string `json:"productId"`
	AddressID string `json:"addressId"`
}

type CheckoutResponse struct {
	Product    []*Product              `json:"product,omitempty"`
	TotalPrice float64                 `json:"totalPrice,omitempty"`
	Cart       *CartResponse           `json:"cart,omitempty"`
	Address    *DefaultAddressResponse `json:"address"`
}
