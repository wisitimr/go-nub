package document

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
