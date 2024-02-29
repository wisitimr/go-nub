package repository

import (
	"context"
	mCollection "findigitalservice/http/rest/internal/model/collection"
	mPaymentMethod "findigitalservice/http/rest/internal/model/payment_method"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	"findigitalservice/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type paymentMethodRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitPaymentMethodRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.PaymentMethodRepository {
	return &paymentMethodRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r paymentMethodRepository) FindAll(ctx context.Context, query map[string][]string) ([]mPaymentMethod.PaymentMethod, error) {
	paymentMethods := []mPaymentMethod.PaymentMethod{}
	cur, err := r.Collection.PaymentMethod.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return paymentMethods, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single paymentMethod can be decoded
		var e mPaymentMethod.PaymentMethod
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		paymentMethods = append(paymentMethods, e)
	}
	return paymentMethods, nil
}

func (r paymentMethodRepository) FindById(ctx context.Context, id string) (mPaymentMethod.PaymentMethod, error) {
	paymentMethod := mPaymentMethod.PaymentMethod{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return paymentMethod, err
	}
	err = r.Collection.PaymentMethod.FindOne(ctx, bson.M{"_id": doc}).Decode(&paymentMethod)
	if err != nil {
		return paymentMethod, err
	}
	return paymentMethod, nil
}

func (r paymentMethodRepository) Create(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error) {
	if _, err := r.Collection.PaymentMethod.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r paymentMethodRepository) Update(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error) {
	var updated mPaymentMethod.PaymentMethod
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"name":      payload.Name,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.PaymentMethod.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
