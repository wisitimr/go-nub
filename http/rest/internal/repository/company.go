package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mCompany "saved/http/rest/internal/model/company"
	mRepo "saved/http/rest/internal/model/repository"
	"saved/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type companyRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitCompanyRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.CompanyRepository {
	return &companyRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r companyRepository) FindAll(ctx context.Context, query map[string][]string) ([]mCompany.Company, error) {
	companys := []mCompany.Company{}
	cur, err := r.Collection.Company.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return companys, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mCompany.Company
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		companys = append(companys, e)
	}
	return companys, nil
}

func (r companyRepository) FindById(ctx context.Context, id string) (mCompany.Company, error) {
	company := mCompany.Company{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return company, err
	}
	err = r.Collection.Company.FindOne(ctx, bson.M{"_id": doc}).Decode(&company)
	if err != nil {
		return company, err
	}
	return company, nil
}

func (r companyRepository) Create(ctx context.Context, payload mCompany.Company) (mCompany.Company, error) {
	if _, err := r.Collection.Company.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r companyRepository) Update(ctx context.Context, payload mCompany.Company) (mCompany.Company, error) {
	var updated mCompany.Company
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":        payload.Code,
			"name":        payload.Name,
			"description": payload.Description,
			"address":     payload.Address,
			"phone":       payload.Phone,
			"contact":     payload.Contact,
			"updatedBy":   payload.UpdatedBy,
			"updatedAt":   payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.Company.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
