package services

import (
	"backend/models"
	"backend/repositories"
	user "backend/services/user_service"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetProductByID(ctx context.Context, productID bson.ObjectID) (*models.Product, error) {

	return repositories.GetProductByID(ctx, productID)
}

func GetSingleCategoryProduct(ctx context.Context) ([]bson.M, error) {

	return repositories.GetSingleCategoryProduct(ctx)
}

func GetCategoryProducts(ctx context.Context, category string) ([]bson.M, error) {

	return repositories.GetCategoryProducts(ctx, category)
}

func GetProducts(ctx context.Context) ([]models.Product, error) {

	return repositories.GetProducts(ctx)
}

func UpdateProduct(ctx context.Context, userID string, productId string, updateData map[string]interface{}) (*models.Product, error) {

	user, err := user.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != "ADMIN" {
		return nil, errors.New("unauthorized")
	}

	objID, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return nil, err
	}

	allowedFields := map[string]bool{
		"productName":  true,
		"brandName":    true,
		"category":     true,
		"productImage": true,
		"description":  true,
		"price":        true,
		"sellingPrice": true,
	}

	safeUpdateData := bson.M{}
	invalidFields := []string{}

	for key, value := range updateData {
		if allowedFields[key] {
			safeUpdateData[key] = value
		} else {
			invalidFields = append(invalidFields, key)
		}
	}

	if len(invalidFields) > 0 {
		return nil, fmt.Errorf("invalid fields")
	}

	if len(safeUpdateData) == 0 {
		return nil, errors.New("no valid fields provided")
	}

	return repositories.UpdateProduct(ctx, objID, safeUpdateData)
}

func UploadProduct(ctx context.Context, userID string, productData *models.Product) (*models.Product, error) {

	user, err := user.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != "ADMIN" {
		return nil, errors.New("unauthorized")
	}

	var productValidation = productData.ProductName == "" || productData.BrandName == "" || productData.Category == "" || len(productData.ProductImage) == 0 || productData.Description == "" || productData.Price <= 0 || productData.SellingPrice <= 0

	if productValidation {
		return nil, errors.New("fields cannot be empty")
	}

	return repositories.CreateProduct(ctx, productData)
}
