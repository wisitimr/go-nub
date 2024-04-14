package repository

import (
	"context"
	mCollection "nub/internal/model/collection"
	mDaybookDetail "nub/internal/model/daybook_detail"
	mRepo "nub/internal/model/repository"
	"nub/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daybookDetailRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitDaybookDetailRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.DaybookDetailRepository {
	return &daybookDetailRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r daybookDetailRepository) FindAll(ctx context.Context, query map[string][]string) ([]mDaybookDetail.DaybookDetail, error) {
	daybookDetails := []mDaybookDetail.DaybookDetail{}
	cur, err := r.Collection.DaybookDetail.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return daybookDetails, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mDaybookDetail.DaybookDetail
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		daybookDetails = append(daybookDetails, e)
	}
	return daybookDetails, nil
}

func (r daybookDetailRepository) FindById(ctx context.Context, id string) (mDaybookDetail.DaybookDetail, error) {
	daybookDetail := mDaybookDetail.DaybookDetail{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return daybookDetail, err
	}
	err = r.Collection.DaybookDetail.FindOne(ctx, bson.M{"_id": doc}).Decode(&daybookDetail)
	if err != nil {
		return daybookDetail, err
	}
	return daybookDetail, nil
}

func (r daybookDetailRepository) Create(ctx context.Context, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error) {
	if _, err := r.Collection.DaybookDetail.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r daybookDetailRepository) Update(ctx context.Context, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error) {
	var updated mDaybookDetail.DaybookDetail
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"name":      payload.Name,
			"detail":    payload.Detail,
			"type":      payload.Type,
			"amount":    payload.Amount,
			"account":   payload.Account,
			"daybook":   payload.Daybook,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.DaybookDetail.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
