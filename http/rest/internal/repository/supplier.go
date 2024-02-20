package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mRepo "saved/http/rest/internal/model/repository"
	mSupplier "saved/http/rest/internal/model/supplier"
	"saved/http/rest/internal/util"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type supplierRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitSupplierRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.SupplierRepository {
	return &supplierRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r supplierRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.Supplier.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r supplierRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.Supplier.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r supplierRepository) FindAll(ctx context.Context, query map[string][]string) ([]mSupplier.Supplier, error) {
	suppliers := []mSupplier.Supplier{}
	cur, err := r.Collection.Supplier.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return suppliers, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mSupplier.Supplier
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		suppliers = append(suppliers, e)
	}
	sort.Slice(suppliers[:], func(i, j int) bool {
		return suppliers[i].Code < suppliers[j].Code
	})
	return suppliers, nil
}

func (r supplierRepository) FindById(ctx context.Context, id string) (mSupplier.Supplier, error) {
	supplier := mSupplier.Supplier{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return supplier, err
	}
	err = r.Collection.Supplier.FindOne(ctx, bson.M{"_id": doc}).Decode(&supplier)
	if err != nil {
		return supplier, err
	}
	return supplier, nil
}

func (r supplierRepository) Create(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error) {
	if _, err := r.Collection.Supplier.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r supplierRepository) Update(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error) {
	var updated mSupplier.Supplier
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":      payload.Code,
			"name":      payload.Name,
			"address":   payload.Address,
			"tax":       payload.Tax,
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

	err := r.Collection.Supplier.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
