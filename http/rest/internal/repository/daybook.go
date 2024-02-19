package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mDaybook "saved/http/rest/internal/model/daybook"
	mRepo "saved/http/rest/internal/model/repository"
	"saved/http/rest/internal/util"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type daybookRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitDaybookRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.DaybookRepository {
	return &daybookRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r daybookRepository) Count(ctx context.Context, query map[string][]string) (int64, error) {
	count, err := r.Collection.Daybook.CountDocuments(ctx, util.QueryHandler(query))
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r daybookRepository) FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error) {
	daybooks := []mDaybook.DaybookList{}
	pipeline := []bson.M{
		{
			"$match": util.QueryHandler(query),
		},
		{
			"$lookup": bson.M{
				"from":         "documents",
				"localField":   "document",
				"foreignField": "_id",
				"as":           "document",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$document",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "companies",
				"localField":   "company",
				"foreignField": "_id",
				"as":           "company",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$company",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$sort": bson.M{
				"createdAt": 1,
			},
		},
	}

	cur, err := r.Collection.Daybook.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error(err)
	}
	if err = cur.All(ctx, &daybooks); err != nil {
		r.logger.Error(err)
	}
	return daybooks, nil
}

func (r daybookRepository) FindAllDetail(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookResponse, error) {
	daybooks := []mDaybook.DaybookResponse{}
	pipeline := []bson.M{
		{
			"$match": util.QueryHandler(query),
		},
		{
			"$lookup": bson.M{
				"from":         "companies",
				"localField":   "company",
				"foreignField": "_id",
				"as":           "company",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$company",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "partners",
				"localField":   "partner",
				"foreignField": "_id",
				"as":           "partner",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$partner",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				// Define the details collection for the join.
				"from": "daybook_details",
				// Specify the variable to use in the pipeline stage.
				"let": bson.M{
					"daybookDetails": "$daybookDetails",
				},
				"pipeline": []bson.M{
					// Select only the relevant details from the details collection.
					// Otherwise all the details are selected.
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$_id",
									"$$daybookDetails",
								},
							},
						},
					},
					// Sort details by their createdAt field in asc. -1 = desc
					{
						"$sort": bson.M{
							"createdAt": 1,
						},
					},
					{
						"$lookup": bson.M{
							"from":         "accounts",
							"localField":   "account",
							"foreignField": "_id",
							"as":           "account",
						},
					},
					{
						"$unwind": bson.M{
							"path":                       "$account",
							"preserveNullAndEmptyArrays": true,
						},
					},
				},
				// Use details as the field name to match struct field.
				"as": "daybookDetails",
			},
		},
	}

	cur, err := r.Collection.Daybook.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error(err)
	}
	if err = cur.All(ctx, &daybooks); err != nil {
		r.logger.Error(err)
	}
	for i, inv := range daybooks {
		for _, doc := range inv.DaybookDetails {
			switch doc.Type {
			case "DR":
				daybooks[i].DebitTotalAmount += doc.Amount
			case "CR":
				daybooks[i].CreditTotalAmount += doc.Amount
			}
		}
	}
	return daybooks, nil
}

func (r daybookRepository) FindById(ctx context.Context, id string) (mDaybook.DaybookResponse, error) {
	daybook := mDaybook.Daybook{}
	result := mDaybook.DaybookResponse{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDaybook.DaybookResponse{}, err
	}
	err = r.Collection.Daybook.FindOne(ctx, bson.M{"_id": doc}).Decode(&daybook)
	if err != nil {
		return result, err
	}
	r.logger.Info(daybook)
	result.Id = daybook.Id
	result.Number = daybook.Number
	result.Invoice = daybook.Invoice
	result.Document = daybook.Document
	result.TransactionDate = daybook.TransactionDate
	result.Company = daybook.Company
	result.Supplier = daybook.Supplier
	result.Customer = daybook.Customer
	result.CreatedBy = daybook.CreatedBy
	result.CreatedAt = daybook.CreatedAt
	result.UpdatedBy = daybook.UpdatedBy
	result.UpdatedAt = daybook.UpdatedAt
	outDaybookDetails := []mDaybook.OutDaybookDetails{}
	daybookDetails := []mDaybook.DaybookDetails{}
	ch := make(chan mDaybook.OutDaybookDetails)
	for _, doc := range daybook.DaybookDetails {
		go func(ch chan mDaybook.OutDaybookDetails, doc primitive.ObjectID) {
			var out []mDaybook.OutDaybookDetails
			// err = r.Collection.DaybookDetail.FindOne(ctx, bson.M{"_id": doc}).Decode(&out)
			pipeline := []bson.M{
				{
					"$match": bson.M{"_id": doc},
				},
				{
					"$lookup": bson.M{
						"from":         "accounts",
						"localField":   "account",
						"foreignField": "_id",
						"as":           "account",
					},
				},
				{
					"$unwind": bson.M{
						"path":                       "$account",
						"preserveNullAndEmptyArrays": true,
					},
				},
			}
			cur, err := r.Collection.DaybookDetail.Aggregate(ctx, pipeline)
			if err != nil {
				r.logger.Error(err)
				return
			}
			if err = cur.All(ctx, &out); err != nil {
				r.logger.Error(err)
				return
			}
			ch <- out[0]
		}(ch, doc)
	}
	for range daybook.DaybookDetails {
		outDaybookDetails = append(outDaybookDetails, <-ch)
	}
	sort.Slice(outDaybookDetails[:], func(i, j int) bool {
		return outDaybookDetails[i].CreatedAt.Before(outDaybookDetails[j].CreatedAt)
	})
	for _, row := range outDaybookDetails {
		daybookDetail := mDaybook.DaybookDetails{}
		daybookDetail.Id = row.Id
		daybookDetail.Name = row.Name
		daybookDetail.Type = row.Type
		daybookDetail.Amount = row.Amount
		daybookDetail.Account = row.Account
		daybookDetails = append(daybookDetails, daybookDetail)
		switch row.Type {
		case "DR":
			result.DebitTotalAmount += row.Amount
		case "CR":
			result.CreditTotalAmount += row.Amount
		}
	}
	result.DaybookDetails = daybookDetails
	return result, nil
}

func (r daybookRepository) Create(ctx context.Context, payload mDaybook.Daybook) (mDaybook.Daybook, error) {
	if _, err := r.Collection.Daybook.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r daybookRepository) BulkCreateDaybookDetail(ctx context.Context, payloads []interface{}) error {
	_, err := r.Collection.DaybookDetail.InsertMany(ctx, payloads)
	if err != nil {
		r.logger.Error(err)
		return err
	}
	return nil
}

func (r daybookRepository) Update(ctx context.Context, payload mDaybook.Daybook) (mDaybook.Daybook, error) {
	var updated mDaybook.Daybook
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"number":          payload.Number,
			"invoice":         payload.Invoice,
			"document":        payload.Document,
			"transactionDate": payload.TransactionDate,
			"company":         payload.Company,
			"supplier":        payload.Supplier,
			"customer":        payload.Customer,
			"daybookDetails":  payload.DaybookDetails,
			"updatedBy":       payload.UpdatedBy,
			"updatedAt":       payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := r.Collection.Daybook.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
