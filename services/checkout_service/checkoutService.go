package checkoutservice

import (
	"backend/models"
	addressService "backend/services/address_service"
	cartService "backend/services/cart_service"
	productService "backend/services/product_service"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Checkout(ctx context.Context, req models.CheckoutRequest, userID string, userObjId bson.ObjectID) (*models.CheckoutResponse, error) {

	var response *models.CheckoutResponse

	// ADDRESS (always required)
	addressObjId, err := bson.ObjectIDFromHex(req.AddressID)
	if err != nil {
		return nil, err
	}

	address, err := addressService.GetAddress(ctx, addressObjId, userObjId)
	if err != nil {
		return nil, err
	}

	// address = a
	// return nil

	var ErrCartEmpty = errors.New("cart is empty")

	// PRODUCT OR CART
	if req.ProductID == "" {

		cartResp, err := cartService.GetCheckoutCartItems(ctx, userID, *address)
		if err != nil {
			return nil, err
		}

		if len(cartResp.Items) == 0 {
			return nil, ErrCartEmpty
		}

		response = cartResp

	} else {

		productObjId, err := bson.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return nil, err
		}

		productResp, err := productService.GetCheckoutProduct(ctx, productObjId, req.Quantity, *address)
		if err != nil {
			return nil, err
		}

		response = productResp
	}

	// BUILD RESPONSE

	return response, nil

}
