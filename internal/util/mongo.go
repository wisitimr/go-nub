package util

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func QueryHandler(query map[string][]string) bson.M {
	filter := bson.M{}
	for key, value := range query {
		if key == "id" {
			key = "_id"
		}
		doc, err := primitive.ObjectIDFromHex(value[0])
		if err == nil {
			filter[key] = doc
		} else {
			filter[key] = value[0]
		}
	}
	if len(query) > 1 {
		and := bson.M{}
		for key, value := range filter {
			and[key] = value
		}
		filter = bson.M{
			"$and": bson.A{and},
		}
	}

	return filter
}

func JsonToBson(query map[string][]string) bson.M {
	filter := bson.M{}
	for key, value := range query {
		if key == "id" {
			key = "_id"
		}
		doc, err := primitive.ObjectIDFromHex(value[0])
		if err == nil {
			filter[key] = doc
		} else {
			filter[key] = value[0]
		}
	}
	if len(query) > 1 {
		and := bson.M{}
		for key, value := range filter {
			and[key] = value
		}
		filter = bson.M{
			"$and": bson.A{and},
		}
	}

	return filter
}
