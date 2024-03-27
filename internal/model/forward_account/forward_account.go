package forwardAccount

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForwardAccount struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Account   primitive.ObjectID `bson:"account" json:"account"`
	Type      string             `bson:"type" json:"type"`
	Amount    float64            `bson:"amount" json:"amount"`
	Year      int                `bson:"year" json:"year"`
	Company   primitive.ObjectID `bson:"company" json:"company"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
