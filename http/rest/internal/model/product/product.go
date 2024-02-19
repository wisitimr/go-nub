package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	Detail    string             `bson:"detail" json:"detail"`
	Price     float32            `bson:"price" json:"price"`
	Company   primitive.ObjectID `bson:"company" json:"company"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
