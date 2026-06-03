package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProfilePic string        `bson:"profilePic" json:"profilePic"`
	Username   string        `bson:"username" json:"username"`
	Email      string        `bson:"email" json:"email" binding:"required,email"`
	Password   string        `bson:"password,omitempty" json:"-" binding:"required"`
	Role       string        `bson:"role" json:"role"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time     `bson:"updatedAt" json:"updatedAt"`
}
