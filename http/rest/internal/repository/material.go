package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mMaterial "saved/http/rest/internal/model/material"
	mRepo "saved/http/rest/internal/model/repository"
	"saved/http/rest/internal/util"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type materialRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitMaterialRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.MaterialRepository {
	return &materialRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r materialRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.Material.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r materialRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.Material.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r materialRepository) FindAll(ctx context.Context, query map[string][]string) ([]mMaterial.Material, error) {
	materials := []mMaterial.Material{}
	cur, err := r.Collection.Material.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return materials, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mMaterial.Material
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		materials = append(materials, e)
	}
	sort.Slice(materials[:], func(i, j int) bool {
		return materials[i].Code < materials[j].Code
	})
	return materials, nil
}

func (r materialRepository) FindById(ctx context.Context, id string) (mMaterial.Material, error) {
	material := mMaterial.Material{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return material, err
	}
	err = r.Collection.Material.FindOne(ctx, bson.M{"_id": doc}).Decode(&material)
	if err != nil {
		return material, err
	}
	return material, nil
}

func (r materialRepository) Create(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error) {
	if _, err := r.Collection.Material.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r materialRepository) Update(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error) {
	var updated mMaterial.Material
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":        payload.Code,
			"name":        payload.Name,
			"description": payload.Description,
			"company":     payload.Company,
			"updatedBy":   payload.UpdatedBy,
			"updatedAt":   payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.Material.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
