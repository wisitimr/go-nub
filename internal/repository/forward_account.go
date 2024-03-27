package repository

import (
	"context"
	mCollection "findigitalservice/internal/model/collection"
	mForwardAccount "findigitalservice/internal/model/forward_account"
	mRepo "findigitalservice/internal/model/repository"
	"findigitalservice/internal/util"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ForwardAccountRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitForwardAccountRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.ForwardAccountRepository {
	return &ForwardAccountRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r ForwardAccountRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.ForwardAccount.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r ForwardAccountRepository) FindAll(ctx context.Context, query map[string][]string) ([]mForwardAccount.ForwardAccount, error) {
	ForwardAccounts := []mForwardAccount.ForwardAccount{}
	cur, err := r.Collection.ForwardAccount.Find(ctx, util.QueryHandler(query), options.Find().SetSort(bson.D{{Key: "code", Value: 1}}))
	if err != nil {
		return ForwardAccounts, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mForwardAccount.ForwardAccount
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		ForwardAccounts = append(ForwardAccounts, e)
	}
	return ForwardAccounts, nil
}

func (r ForwardAccountRepository) FindById(ctx context.Context, id string) (mForwardAccount.ForwardAccount, error) {
	forwardAccount := mForwardAccount.ForwardAccount{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return forwardAccount, err
	}
	err = r.Collection.ForwardAccount.FindOne(ctx, bson.M{"_id": doc}).Decode(&forwardAccount)
	if err != nil {
		return forwardAccount, err
	}
	return forwardAccount, nil
}

func (r ForwardAccountRepository) FindOne(ctx context.Context, query map[string][]string) (mForwardAccount.ForwardAccount, error) {
	forwardAccount := mForwardAccount.ForwardAccount{}
	filter := bson.M{}
	for key, value := range query {
		if key == "id" {
			key = "_id"
		}
		doc, err := primitive.ObjectIDFromHex(value[0])
		if err == nil {
			filter[key] = doc
		} else {
			if key == "year" {
				yearInt, err := strconv.Atoi(value[0])
				if err != nil {
					return forwardAccount, err
				}
				filter[key] = yearInt
			} else {
				filter[key] = value[0]
			}
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
	err := r.Collection.ForwardAccount.FindOne(ctx, filter).Decode(&forwardAccount)
	if err != nil {
		return forwardAccount, err
	}
	return forwardAccount, nil
}

func (r ForwardAccountRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.ForwardAccount.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r ForwardAccountRepository) Create(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error) {
	if _, err := r.Collection.ForwardAccount.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r ForwardAccountRepository) Update(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error) {
	var updated mForwardAccount.ForwardAccount
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"account":   payload.Account,
			"type":      payload.Type,
			"amount":    payload.Amount,
			"year":      payload.Year,
			"company":   payload.Company,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.ForwardAccount.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
