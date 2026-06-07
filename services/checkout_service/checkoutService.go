package checkoutservice

import (
	"backend/models"
	addressService "backend/services/address_service"
	cartService "backend/services/cart_service"
	productService "backend/services/product_service"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/sync/errgroup"
)

func CheckoutService(ctx context.Context, req models.CheckoutRequest, userID string, userObjId bson.ObjectID) (*models.CheckoutResponse, error) {

	var (
		product *models.Product
		cart    *models.CartResponse
		address *models.DefaultAddressResponse
	)

	g, ctx := errgroup.WithContext(ctx)

	var ErrCartEmpty = errors.New("cart is empty")

	// PRODUCT OR CART
	if req.ProductID == "" {

		g.Go(func() error {
			cartResp, err := cartService.GetCartItems(ctx, userID)
			if err != nil {
				return err
			}

			if len(cartResp.Items) == 0 {
				return ErrCartEmpty
			}

			cart = cartResp
			return nil
		})

	} else {

		productObjId, err := bson.ObjectIDFromHex(req.ProductID)
		if err != nil {
			return nil, err
		}

		g.Go(func() error {
			p, err := productService.GetProductByID(ctx, productObjId)
			if err != nil {
				return err
			}

			product = p
			return nil
		})
	}

	// ADDRESS (always required)
	addressObjId, err := bson.ObjectIDFromHex(req.AddressID)
	if err != nil {
		return nil, err
	}

	g.Go(func() error {
		a, err := addressService.GetAddress(ctx, addressObjId, userObjId)
		if err != nil {
			return err
		}

		address = a
		return nil
	})

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	// BUILD RESPONSE
	resp := &models.CheckoutResponse{
		Cart:    cart,
		Address: address,
	}

	if product != nil {
		resp.Product = []*models.Product{product}
		resp.TotalPrice = product.SellingPrice
	}

	return resp, nil
}
