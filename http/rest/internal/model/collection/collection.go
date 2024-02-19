package collection

import "go.mongodb.org/mongo-driver/mongo"

type Collection struct {
	User          *mongo.Collection
	Account       *mongo.Collection
	Supplier      *mongo.Collection
	Customer      *mongo.Collection
	Document      *mongo.Collection
	Product       *mongo.Collection
	Company       *mongo.Collection
	Daybook       *mongo.Collection
	DaybookDetail *mongo.Collection
	Role          *mongo.Collection
	Material      *mongo.Collection
}