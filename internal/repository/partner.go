package repository

import (
	"context"
	mPartner "findigitalservice/internal/model/partner"
	mRepo "findigitalservice/internal/model/repository"
	"findigitalservice/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type partnerRepository struct {
	Partner *mongo.Collection
	logger  *logrus.Logger
}

func InitPartnerRepository(db *mongo.Database, logger *logrus.Logger) mRepo.PartnerRepository {
	return &partnerRepository{
		Partner: db.Collection("partners"),
		logger:  logger,
	}
}

func (r partnerRepository) FindAll(ctx context.Context, query map[string][]string) ([]mPartner.Partner, error) {
	partners := []mPartner.Partner{}
	cur, err := r.Partner.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return partners, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mPartner.Partner
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		partners = append(partners, e)
	}
	return partners, nil
}

func (r partnerRepository) FindById(ctx context.Context, id string) (mPartner.Partner, error) {
	partner := mPartner.Partner{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return partner, err
	}
	err = r.Partner.FindOne(ctx, bson.M{"_id": doc}).Decode(&partner)
	if err != nil {
		return partner, err
	}
	return partner, nil
}

func (r partnerRepository) Create(ctx context.Context, payload mPartner.Partner) (mPartner.Partner, error) {
	if _, err := r.Partner.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r partnerRepository) Update(ctx context.Context, payload mPartner.Partner) (mPartner.Partner, error) {
	var updated mPartner.Partner
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":      payload.Code,
			"name":      payload.Name,
			"address":   payload.Address,
			"tax":       payload.Tax,
			"phone":     payload.Phone,
			"contact":   payload.Contact,
			"type":      payload.Type,
			"company":   payload.Company,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := r.Partner.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
