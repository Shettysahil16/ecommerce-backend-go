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

	//var response *models.CheckoutResponse

	// ADDRESS (always required)
	// addressObjId, err := bson.ObjectIDFromHex(req.AddressID)
	// if err != nil {
	// 	return nil, err
	// }

	// address, err := addressService.GetAddress(ctx, addressObjId, userObjId)
	// if err != nil {
	// 	return nil, err
	// }

	// address = a
	// return nil

	// PRODUCT OR CART
	if req.ProductID == "" {

		cartResp, err := cartService.GetCheckoutCartItems(ctx, userID)
		if err != nil {
			return err
		}

		if len(cartResp.Items) == 0 {
			return ErrCartEmpty
		}

		//response = cartResp

		var checkoutItems []models.CheckoutItem

		for _, item := range cartResp.Items {
			checkoutItems = append(checkoutItems, models.CheckoutItem{
				ProductID: item.Product.ID.Hex(),
				Quantity:  int64(item.Quantity),
			})
		}

		err = cache.SetCheckoutItems(ctx, userID, checkoutItems)
		if err != nil {
			return err
		}

	} else {

		productObjId, err := bson.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return err
		}

		_, err = productService.GetCheckoutProduct(ctx, productObjId, req.Quantity)
		if err != nil {
			return err
		}

		//response = productResp

		err = cache.SetCheckoutItems(ctx, userID, []models.CheckoutItem{
			{
				ProductID: req.ProductID,
				Quantity:  int64(req.Quantity),
			},
		})
		if err != nil {
			return err
		}

	}

	// BUILD RESPONSE

	return nil

}

func PrepareCheckout(ctx context.Context, userID string, userObjId bson.ObjectID, req models.CheckoutRequest) (*models.CheckoutResponse, error) {

	checkoutData, err := cache.GetCheckout(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(checkoutData) == 0 {
		if err = SetCheckoutItems(ctx, req, userID); err != nil {
			return nil, err
		}

		checkoutData, err = cache.GetCheckout(ctx, userID)
		if err != nil {
			return nil, err
		}
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

		qty := cartMap[p.ID.Hex()]

		result = append(result, models.CheckoutItemResponse{
			Product:  p,
			Quantity: qty,
		})

		totalPrice += p.SellingPrice * float64(qty)
	}

	// ADDRESS (always required)
	addressObjId, err := bson.ObjectIDFromHex(req.AddressID)
	if err != nil {
		return nil, err
	}

	address, err := addressService.GetAddress(ctx, addressObjId, userObjId)
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
