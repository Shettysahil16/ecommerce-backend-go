package checkoutservice

import (
	"backend/cache"
	"backend/models"
	"backend/repositories"
	addressService "backend/services/address_service"
	cartService "backend/services/cart_service"
	productService "backend/services/product_service"
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrNoCheckoutItems = errors.New("no checkout items found")
var ErrCheckoutProductsUnavailable = errors.New("some products are unavailable")

func SetCheckoutItems(ctx context.Context, req models.CheckoutRequest, userID string) error {

	var checkoutItems []models.CheckoutItem

	if req.ProductID == "" {

		// CART
		cartResp, err := cartService.GetCheckoutCartItems(ctx, userID)
		if err != nil {
			return err
		}

		if len(cartResp.Items) == 0 {
			return ErrCartEmpty
		}

		for _, item := range cartResp.Items {
			checkoutItems = append(checkoutItems, models.CheckoutItem{
				ProductID: item.Product.ID.Hex(),
				Quantity:  int64(item.Quantity),
			})
		}

	} else {

		// PRODUCT
		productObjId, err := bson.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return err
		}

		product, err := productService.GetProductByID(ctx, productObjId)
		if err != nil {
			return err
		}

		qty := req.Quantity
		if qty <= 0 {
			qty = 1
		}

		checkoutItems = []models.CheckoutItem{
			{
				ProductID: product.ID.Hex(),
				Quantity:  int64(qty),
			},
		}
	}

	// BUILD RESPONSE
	return cache.SetCheckoutItems(ctx, userID, checkoutItems)
}

func PrepareCheckout(ctx context.Context, userID string, userObjId bson.ObjectID, addressID bson.ObjectID) (*models.CheckoutResponse, error) {

	checkoutData, err := cache.GetCheckout(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(checkoutData) == 0 {
		return nil, ErrNoCheckoutItems
	}

	var productIDs []bson.ObjectID
	cartMap := make(map[string]int)

	for productID, qty := range checkoutData {
		objID, err := bson.ObjectIDFromHex(productID)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(qty)
		if err != nil {
			return nil, err
		}

		productIDs = append(productIDs, objID)
		cartMap[productID] = quantity
	}

	if len(productIDs) == 0 {
		return nil, ErrNoCheckoutItems
	}

	products, err := repositories.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	if len(products) != len(productIDs) {
		return nil, ErrCheckoutProductsUnavailable
	}

	var result []models.CheckoutItemResponse
	var totalPrice float64

	for _, p := range products {

		qty, ok := cartMap[p.ID.Hex()]
		if !ok {
			qty = 1
		}

		result = append(result, models.CheckoutItemResponse{
			Product:  p,
			Quantity: qty,
		})

		totalPrice += p.SellingPrice * float64(qty)
	}

	// ADDRESS (always required)

	address, err := addressService.GetAddress(ctx, addressID, userObjId)
	if err != nil {
		return nil, err
	}

	//address = a
	//return nil

	return &models.CheckoutResponse{
		Items:      result,
		TotalPrice: totalPrice,
		Address:    address,
	}, nil

}

func UpdateCheckoutItemService(ctx context.Context, userID string, productID string, qty int64) error {

	if productID == "" {
		return errors.New("invalid product id")
	}
	if qty <= 0 {
		return cache.RemoveCheckoutItem(ctx, userID, productID)
	}

	return cache.UpdateCheckoutItem(ctx, userID, productID, qty)
}

func RemoveCheckoutItemService(ctx context.Context, userID string, productID string) error {

	if productID == "" {
		return errors.New("invalid product id")
	}

	return cache.RemoveCheckoutItem(ctx, userID, productID)
}
