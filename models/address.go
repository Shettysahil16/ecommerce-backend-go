package models

import "go.mongodb.org/mongo-driver/v2/bson"

type Address struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID `bson:"userId" json:"userId"`
	FullName    string        `bson:"fullName" json:"fullName"`
	State       string        `bson:"state" json:"state"`
	City        string        `bson:"city" json:"city"`
	Pincode     string        `bson:"pincode" json:"pincode"`
	Phone       string        `bson:"phone" json:"phone"`
	Area        string        `bson:"area" json:"area"`
	Landmark    string        `bson:"landmark" json:"landmark"`
	HouseNo     string        `bson:"houseNo" json:"houseNo"`
	AddressType string        `bson:"addressType" json:"addressType"`
	IsDefault   bool          `bson:"isDefault" json:"isDefault"`
}

type AddressResponse struct {
	FullName    string `bson:"fullName,omitempty" json:"fullName"`
	State       string `bson:"state,omitempty" json:"state"`
	City        string `bson:"city,omitempty" json:"city"`
	Pincode     string `bson:"pincode,omitempty" json:"pincode"`
	Phone       string `bson:"phone,omitempty" json:"phone"`
	Area        string `bson:"area,omitempty" json:"area"`
	Landmark    string `bson:"landmark,omitempty" json:"landmark"`
	HouseNo     string `bson:"houseNo,omitempty" json:"houseNo"`
	AddressType string `bson:"addressType,omitempty" json:"addressType"`
}

type DefaultAddressResponse struct {
	FullName    string `bson:"fullName,omitempty" json:"fullName"`
	State       string `bson:"state,omitempty" json:"state"`
	City        string `bson:"city,omitempty" json:"city"`
	Pincode     string `bson:"pincode,omitempty" json:"pincode"`
	Phone       string `bson:"phone,omitempty" json:"phone"`
	Area        string `bson:"area,omitempty" json:"area"`
	Landmark    string `bson:"landmark,omitempty" json:"landmark"`
	HouseNo     string `bson:"houseNo,omitempty" json:"houseNo"`
	AddressType string `bson:"addressType,omitempty" json:"addressType"`
	IsDefault   bool   `bson:"isDefault" json:"isDefault"`
}
