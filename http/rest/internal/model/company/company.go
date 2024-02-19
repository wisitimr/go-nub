package company

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Code        string             `bson:"code" json:"code"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Address     string             `bson:"address" json:"address"`
	Phone       string             `bson:"phone" json:"phone"`
	Contact     string             `bson:"contact" json:"contact"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}
