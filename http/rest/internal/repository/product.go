package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mProduct "saved/http/rest/internal/model/product"
	mRepo "saved/http/rest/internal/model/repository"
	"saved/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitProductRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.ProductRepository {
	return &productRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r productRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.Product.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r productRepository) FindAll(ctx context.Context, query map[string][]string) ([]mProduct.Product, error) {
	products := []mProduct.Product{}
	cur, err := r.Collection.Product.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return products, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mProduct.Product
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		products = append(products, e)
	}
	return products, nil
}

func (r productRepository) FindById(ctx context.Context, id string) (mProduct.Product, error) {
	product := mProduct.Product{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return product, err
	}
	err = r.Collection.Product.FindOne(ctx, bson.M{"_id": doc}).Decode(&product)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r productRepository) Create(ctx context.Context, payload mProduct.Product) (mProduct.Product, error) {
	if _, err := r.Collection.Product.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r productRepository) Update(ctx context.Context, payload mProduct.Product) (mProduct.Product, error) {
	var updated mProduct.Product
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":        payload.Code,
			"name":        payload.Name,
			"description": payload.Description,
			"price":       payload.Price,
			"company":     payload.Company,
			"updatedBy":   payload.UpdatedBy,
			"updatedAt":   payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := r.Collection.Product.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
