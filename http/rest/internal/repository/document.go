package repository

import (
	"context"
	mCollection "findigitalservice/http/rest/internal/model/collection"
	mDocument "findigitalservice/http/rest/internal/model/document"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	"findigitalservice/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type documentRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitDocumentRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.DocumentRepository {
	return &documentRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r documentRepository) FindAll(ctx context.Context, query map[string][]string) ([]mDocument.Document, error) {
	documents := []mDocument.Document{}
	cur, err := r.Collection.Document.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return documents, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mDocument.Document
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		documents = append(documents, e)
	}
	return documents, nil
}

func (r documentRepository) FindById(ctx context.Context, id string) (mDocument.Document, error) {
	document := mDocument.Document{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return document, err
	}
	err = r.Collection.Document.FindOne(ctx, bson.M{"_id": doc}).Decode(&document)
	if err != nil {
		return document, err
	}
	return document, nil
}

func (r documentRepository) Create(ctx context.Context, payload mDocument.Document) (mDocument.Document, error) {
	if _, err := r.Collection.Document.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r documentRepository) Update(ctx context.Context, payload mDocument.Document) (mDocument.Document, error) {
	var updated mDocument.Document
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"code":      payload.Code,
			"name":      payload.Name,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	err := r.Collection.Document.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
