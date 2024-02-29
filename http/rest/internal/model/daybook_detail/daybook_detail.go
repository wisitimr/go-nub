package daybookDetail

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DaybookDetail struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Amount    float64            `bson:"amount" json:"amount"`
	Account   primitive.ObjectID `bson:"account" json:"account"`
	Daybook   primitive.ObjectID `bson:"daybook" json:"daybook"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
