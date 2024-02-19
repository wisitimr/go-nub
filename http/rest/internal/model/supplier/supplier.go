package supplier

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Supplier struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	Address   string             `bson:"address" json:"address"`
	Tax       string             `bson:"tax" json:"tax"`
	Phone     string             `bson:"phone" json:"phone"`
	Contact   string             `bson:"contact" json:"contact"`
	Company   primitive.ObjectID `bson:"company" json:"company"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
