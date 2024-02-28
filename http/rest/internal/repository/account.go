package repository

import (
	"context"
	mAccount "findigitalservice/http/rest/internal/model/account"
	mCollection "findigitalservice/http/rest/internal/model/collection"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	"findigitalservice/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type accountRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitAccountRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.AccountRepository {
	return &accountRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r accountRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.Account.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r accountRepository) FindAll(ctx context.Context, query map[string][]string) ([]mAccount.Account, error) {
	accounts := []mAccount.Account{}
	cur, err := r.Collection.Account.Find(ctx, util.QueryHandler(query), options.Find().SetSort(bson.D{{Key: "code", Value: 1}}))
	if err != nil {
		return accounts, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mAccount.Account
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		accounts = append(accounts, e)
	}
	return accounts, nil
}

func (r accountRepository) FindById(ctx context.Context, id string) (mAccount.Account, error) {
	account := mAccount.Account{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return account, err
	}
	err = r.Collection.Account.FindOne(ctx, bson.M{"_id": doc}).Decode(&account)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (r accountRepository) Delete(ctx context.Context, id string) error {
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.Collection.Account.DeleteOne(ctx, bson.M{"_id": doc})
	if err != nil {
		return err
	}
	return nil
}

func (r accountRepository) Create(ctx context.Context, payload mAccount.Account) (mAccount.Account, error) {
	if _, err := r.Collection.Account.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r accountRepository) Update(ctx context.Context, payload mAccount.Account) (mAccount.Account, error) {
	var updated mAccount.Account
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":        payload.Code,
			"name":        payload.Name,
			"description": payload.Description,
			"type":        payload.Type,
			"company":     payload.Company,
			"updatedBy":   payload.UpdatedBy,
			"updatedAt":   payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.Account.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
