package cart

import (
	"backend/cache"
	"backend/models"
	"backend/repositories"
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetCartItems(ctx context.Context, userID string) (*models.CartResponse, error) {

	cartData, err := cache.GetCartCache(userID)
	if err != nil {
		return nil, err
	}

	var productIDs []bson.ObjectID
	cartMap := make(map[string]int)

	for productID, qty := range cartData {
		objID, err := bson.ObjectIDFromHex(productID)
		if err != nil {
			continue
		}

		quantity, _ := strconv.Atoi(qty)

		productIDs = append(productIDs, objID)
		cartMap[productID] = quantity
	}

	if len(productIDs) == 0 {
		return &models.CartResponse{}, nil
	}

	products, err := repositories.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	var result []models.CartItemResponse
	var totalPrice float64

	for _, p := range products {

		qty := cartMap[p.ID.Hex()]

		result = append(result, models.CartItemResponse{
			Product:  p,
			Quantity: qty,
		})

		totalPrice += p.SellingPrice * float64(qty)
	}

	return &models.CartResponse{
		Items:      result,
		TotalPrice: totalPrice,
	}, nil
}

func AddToCart(ctx context.Context, userID string, productID bson.ObjectID, qty int) (*models.RedisCart, error) {

	return cache.AddToCartCache(userID, productID.Hex(), qty)
}

func UpdateCheckoutCartQuantity(ctx context.Context, userID string, productID bson.ObjectID, qty int) (*models.RedisCart, error) {

	if qty < 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	return cache.UpdateCartQuantityCache(userID, productID.Hex(), qty)
}

func DeleteCartItem(ctx context.Context, userID string, productID bson.ObjectID) error {

	return cache.DeleteCartItemCache(userID, productID.Hex())
}

func ReduceCartItem(ctx context.Context, userID string, productID bson.ObjectID, qty int) (*models.RedisCart, error) {

	return cache.ReduceCartCache(userID, productID.Hex(), qty)
}

func GetCheckoutCartItems(ctx context.Context, userID string, address models.DefaultAddressResponse) (*models.CheckoutResponse, error) {

	cartData, err := cache.GetCartCache(userID)
	if err != nil {
		return nil, err
	}

	var productIDs []bson.ObjectID
	cartMap := make(map[string]int)

	for productID, qty := range cartData {
		objID, err := bson.ObjectIDFromHex(productID)
		if err != nil {
			continue
		}

		quantity, err := strconv.Atoi(qty)
		if err != nil {
			return nil, err
		}

		productIDs = append(productIDs, objID)
		cartMap[productID] = quantity
	}

	if len(productIDs) == 0 {
		return &models.CheckoutResponse{}, nil
	}

	products, err := repositories.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	var result []models.CheckoutItemResponse
	var totalPrice float64

	for _, p := range products {

		qty := cartMap[p.ID.Hex()]

		result = append(result, models.CheckoutItemResponse{
			Product:  p,
			Quantity: qty,
		})

		totalPrice += p.SellingPrice * float64(qty)
	}

	return &models.CheckoutResponse{
		Items:      result,
		TotalPrice: totalPrice,
		Address:    &address,
	}, nil
}
