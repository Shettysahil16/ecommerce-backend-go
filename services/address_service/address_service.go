package address

import (
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func isEmpty(fields ...string) bool {

	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return true
		}
	}
	return false
}

func CreateAddress(ctx context.Context, AddressData models.Address, userID bson.ObjectID) (*models.Address, error) {

	// EMPTY FIELDS CHECKING LOGIC
	if isEmpty(
		AddressData.FullName,
		AddressData.State,
		AddressData.City,
		AddressData.Pincode,
		AddressData.Phone,
		AddressData.Area,
		AddressData.Landmark,
		AddressData.HouseNo,
	) {
		return nil, errors.New("fields cannot be empty")
	}

	// PINCODE CHECKING LOGIC
	var pincodeRegex = regexp.MustCompile(`^[1-9][0-9]{5}$`)

	if !pincodeRegex.MatchString(AddressData.Pincode) {
		return nil, errors.New("invalid pincode")
	}

	//api for checking the pincode matching district https://api.postalpincode.in/pincode/400001
	err := utils.ValidatePincodeandLocation(AddressData.Pincode, AddressData.State, AddressData.City)
	if err != nil {
		return nil, err
	}

	// MOBILE NUMBER CHECKING LOGIC
	if len(AddressData.Phone) != 10 {
		return nil, errors.New("invalid mobile number")
	}

	_, err = strconv.Atoi(AddressData.Phone)
	if err != nil {
		return nil, errors.New("mobile number must contain only digits")
	}

	// ADDRESS TYPE MUST BE WORK, HOME OR OTHER
	AddressData.AddressType = strings.ToLower(strings.TrimSpace(AddressData.AddressType))

	if AddressData.AddressType == "" {
		AddressData.AddressType = "home"
	}

	switch AddressData.AddressType {
	case "home", "work", "other":

	default:
		return nil, errors.New("address type must be home, work, or other")
	}

	// SESSION START
	session, err := repositories.AddressCollection().Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// TRANSACTION STARTED
	result, err := session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		totalAddresses, err := CountAddresses(ctx, userID)
		if err != nil {
			return nil, err
		}

		if totalAddresses >= 5 {
			return nil, errors.New("maximum address limit reached")
		}

		if totalAddresses == 0 {
			AddressData.IsDefault = true
		}

		createdAddr, err := repositories.CreateAddress(ctx, AddressData, userID)
		if err != nil {
			return nil, err
		}

		if AddressData.IsDefault {
			err := UpdateManyAddress(ctx, createdAddr.ID, userID)
			if err != nil {
				return nil, err
			}
		}

		return createdAddr, nil
	})

	if err != nil {
		return nil, err
	}

	createdAddress, ok := result.(*models.Address)
	if !ok {
		return nil, errors.New("unexpected transaction result type")
	}

	return createdAddress, nil

}

func GetAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) (*models.DefaultAddressResponse, error) {

	return repositories.GetAddress(ctx, addressID, userID)
}

func GetAddresses(ctx context.Context, userID bson.ObjectID) (*[]models.Address, error) {

	return repositories.GetAddresses(ctx, userID)
}

func UpdateAddress(ctx context.Context, addressID bson.ObjectID, addressData models.AddressResponse) (*models.AddressResponse, error) {

	return repositories.UpdateAddress(ctx, addressID, addressData)
}

func UpdateManyAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) error {

	return repositories.UpdateManyAddress(ctx, addressID, userID)
}

func UpdateDefaultAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) (*models.DefaultAddressResponse, error) {

	err := UpdateManyAddress(ctx, addressID, userID)
	if err != nil {
		return nil, err
	}

	return repositories.UpdateDefaultAddress(ctx, addressID, userID)
}

func GetFirstAddressUpdate(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) (*models.Address, error) {

	return repositories.GetFirstAddressUpdate(ctx, userID, addressID)
}

func DeleteAddress(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) error {

	session, err := repositories.AddressCollection().Database().Client().StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (any, error) {
		address, err := repositories.GetAddress(sc, addressID, userID)
		if err != nil {
			return nil, err
		}

		if address.IsDefault {
			_, err := GetFirstAddressUpdate(sc, userID, addressID)
			if err != nil {
				return nil, err
			}
		}

		if err := repositories.DeleteAddress(sc, userID, addressID); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err

}

func CountAddresses(ctx context.Context, userID bson.ObjectID) (int64, error) {

	return repositories.CountAddresses(ctx, userID)
}

func GetFirstAddress(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) (*models.Address, error) {

	return repositories.GetFirstAddress(ctx, userID, addressID)
}
