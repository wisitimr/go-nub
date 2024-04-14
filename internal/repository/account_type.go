package repository

import (
	"context"
	mAccountType "nub/internal/model/account_type"
	mCollection "nub/internal/model/collection"
	mRepo "nub/internal/model/repository"
	"nub/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountTypeRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitAccountTypeRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.AccountTypeRepository {
	return &accountTypeRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r accountTypeRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.AccountType.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r accountTypeRepository) FindAll(ctx context.Context, query map[string][]string) ([]mAccountType.AccountType, error) {
	accounts := []mAccountType.AccountType{}
	cur, err := r.Collection.AccountType.Find(ctx, util.QueryHandler(query), options.Find().SetSort(bson.D{{Key: "code", Value: 1}}))
	if err != nil {
		return accounts, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mAccountType.AccountType
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		accounts = append(accounts, e)
	}
	return accounts, nil
}

func (r accountTypeRepository) FindById(ctx context.Context, id string) (mAccountType.AccountType, error) {
	account := mAccountType.AccountType{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return account, err
	}
	err = r.Collection.AccountType.FindOne(ctx, bson.M{"_id": doc}).Decode(&account)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (r accountTypeRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.AccountType.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r accountTypeRepository) Create(ctx context.Context, payload mAccountType.AccountType) (mAccountType.AccountType, error) {
	if _, err := r.Collection.AccountType.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r accountTypeRepository) Update(ctx context.Context, payload mAccountType.AccountType) (mAccountType.AccountType, error) {
	var updated mAccountType.AccountType
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"name":      payload.Name,
			"company":   payload.Company,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.AccountType.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
