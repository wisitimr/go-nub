package repository

import (
	"context"
	mCollection "findigitalservice/http/rest/internal/model/collection"
	mCustomer "findigitalservice/http/rest/internal/model/customer"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	"findigitalservice/http/rest/internal/util"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type customerRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitCustomerRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.CustomerRepository {
	return &customerRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r customerRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.Customer.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r customerRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.Customer.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r customerRepository) FindAll(ctx context.Context, query map[string][]string) ([]mCustomer.Customer, error) {
	customers := []mCustomer.Customer{}
	cur, err := r.Collection.Customer.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return customers, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mCustomer.Customer
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		customers = append(customers, e)
	}
	sort.Slice(customers[:], func(i, j int) bool {
		return customers[i].Code < customers[j].Code
	})
	return customers, nil
}

func (r customerRepository) FindById(ctx context.Context, id string) (mCustomer.Customer, error) {
	customer := mCustomer.Customer{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customer, err
	}
	err = r.Collection.Customer.FindOne(ctx, bson.M{"_id": doc}).Decode(&customer)
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func (r customerRepository) Create(ctx context.Context, payload mCustomer.Customer) (mCustomer.Customer, error) {
	if _, err := r.Collection.Customer.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r customerRepository) Update(ctx context.Context, payload mCustomer.Customer) (mCustomer.Customer, error) {
	var updated mCustomer.Customer
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":      payload.Code,
			"name":      payload.Name,
			"address":   payload.Address,
			"phone":     payload.Phone,
			"contact":   payload.Contact,
			"company":   payload.Company,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := r.Collection.Customer.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
