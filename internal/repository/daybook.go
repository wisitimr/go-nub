package repository

import (
	"context"
	"fmt"
	mCollection "nub/internal/model/collection"
	mDaybook "nub/internal/model/daybook"
	mRepo "nub/internal/model/repository"
	"sort"
	"strconv"
	"time"

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
	filter := bson.M{}
	if query["company"] != nil {
		doc, _ := primitive.ObjectIDFromHex(query["company"][0])
		filter["company"] = bson.M{"$eq": doc}
	}
	const (
		layoutISO = "2006-01-02T15:04:05.000Z"
	)
	if query["transactionDate.gte"] != nil && query["transactionDate.lt"] != nil {
		from, _ := time.Parse(layoutISO, query["transactionDate.gte"][0])
		to, _ := time.Parse(layoutISO, query["transactionDate.lt"][0])
		filter["transactionDate"] = bson.M{
			"$gte": from,
			"$lt":  to,
		}
	} else if query["transactionDate.gte"] != nil {
		from, _ := time.Parse(layoutISO, query["transactionDate.gte"][0])
		filter["transactionDate"] = bson.M{
			"$gte": from,
		}
	} else if query["transactionDate.lt"] != nil {
		to, _ := time.Parse(layoutISO, query["transactionDate.lt"][0])
		filter["transactionDate"] = bson.M{
			"$lt": to,
		}
	}
	if len(filter) > 1 {
		and := bson.M{}
		for key, value := range filter {
			and[key] = value
		}
		filter = bson.M{
			"$and": bson.A{and},
		}
	}
	count, err := r.Collection.Daybook.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r daybookRepository) FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error) {
	daybooks := []mDaybook.DaybookList{}
	filter := bson.M{}
	if query["company"] != nil {
		doc, _ := primitive.ObjectIDFromHex(query["company"][0])
		filter["company"] = bson.M{"$eq": doc}
	}
	const (
		layoutISO = "2006-01-02T15:04:05.000Z"
	)
	if query["transactionDate"] != nil {
		d, _ := time.Parse(layoutISO, query["transactionDate"][0])
		filter["transactionDate"] = d
	} else {
		if query["transactionDate.gte"] != nil && query["transactionDate.lte"] != nil {
			from, _ := time.Parse(layoutISO, query["transactionDate.gte"][0])
			to, _ := time.Parse(layoutISO, query["transactionDate.lte"][0])
			filter["transactionDate"] = bson.M{
				"$gte": from,
				"$lte": to,
			}
		} else if query["transactionDate.gte"] != nil && query["transactionDate.lt"] != nil {
			from, _ := time.Parse(layoutISO, query["transactionDate.gte"][0])
			to, _ := time.Parse(layoutISO, query["transactionDate.lt"][0])
			filter["transactionDate"] = bson.M{
				"$gte": from,
				"$lt":  to,
			}
		} else if query["transactionDate.gt"] != nil && query["transactionDate.lte"] != nil {
			from, _ := time.Parse(layoutISO, query["transactionDate.gt"][0])
			to, _ := time.Parse(layoutISO, query["transactionDate.lte"][0])
			filter["transactionDate"] = bson.M{
				"$gt":  from,
				"$lte": to,
			}
		} else if query["transactionDate.gte"] != nil {
			from, _ := time.Parse(layoutISO, query["transactionDate.gte"][0])
			filter["transactionDate"] = bson.M{
				"$gte": from,
			}
		} else if query["transactionDate.lte"] != nil {
			to, _ := time.Parse(layoutISO, query["transactionDate.lte"][0])
			filter["transactionDate"] = bson.M{
				"$lte": to,
			}
		}
	}
	if len(filter) > 1 {
		and := bson.M{}
		for key, value := range filter {
			and[key] = value
		}
		filter = bson.M{
			"$and": bson.A{and},
		}
	}
	pipeline := []bson.M{
		{
			"$match": filter,
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
			"$lookup": bson.M{
				"from":         "suppliers",
				"localField":   "supplier",
				"foreignField": "_id",
				"as":           "supplier",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$supplier",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "customers",
				"localField":   "customer",
				"foreignField": "_id",
				"as":           "customer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$customer",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$sort": bson.M{
				"transactionDate": -1,
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
	result.Id = daybook.Id
	result.Number = daybook.Number
	result.Invoice = daybook.Invoice
	result.Document = daybook.Document
	result.TransactionDate = daybook.TransactionDate
	result.Company = daybook.Company
	result.Supplier = daybook.Supplier
	result.Customer = daybook.Customer
	result.PaymentMethod = daybook.PaymentMethod
	result.CreatedBy = daybook.CreatedBy
	result.CreatedAt = daybook.CreatedAt
	result.UpdatedBy = daybook.UpdatedBy
	result.UpdatedAt = daybook.UpdatedAt
	outDaybookDetails := []mDaybook.OutDaybookDetails{}
	daybookDetails := []mDaybook.DaybookDetail{}
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
		daybookDetail := mDaybook.DaybookDetail{}
		daybookDetail.Id = row.Id
		daybookDetail.Name = row.Name
		daybookDetail.Type = row.Type
		daybookDetail.Amount = row.Amount
		daybookDetail.Account = row.Account
		daybookDetail.CreatedBy = row.CreatedBy
		daybookDetail.CreatedAt = row.CreatedAt
		daybookDetail.UpdatedBy = row.UpdatedBy
		daybookDetail.UpdatedAt = row.UpdatedAt
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

func (r daybookRepository) FindByIdForExcel(ctx context.Context, id string) (mDaybook.DaybookExpand, error) {
	daybooks := []mDaybook.DaybookExpand{}
	// result := mDaybook.DaybookResponse{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDaybook.DaybookExpand{}, err
	}
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": doc},
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
			"$lookup": bson.M{
				"from":         "suppliers",
				"localField":   "supplier",
				"foreignField": "_id",
				"as":           "supplier",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$supplier",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "customers",
				"localField":   "customer",
				"foreignField": "_id",
				"as":           "customer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$customer",
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
	daybook := daybooks[0]
	sort.Slice(daybook.DaybookDetails[:], func(i, j int) bool {
		return daybook.DaybookDetails[i].CreatedAt.Before(daybook.DaybookDetails[j].CreatedAt)
	})
	daybookDetails := []mDaybook.DaybookDetail{}
	for _, row := range daybook.DaybookDetails {
		daybookDetail := mDaybook.DaybookDetail{}
		daybookDetail.Id = row.Id
		daybookDetail.Name = row.Name
		daybookDetail.Type = row.Type
		daybookDetail.Amount = row.Amount
		daybookDetail.Account = row.Account
		daybookDetails = append(daybookDetails, daybookDetail)
		switch row.Type {
		case "DR":
			daybook.DebitTotalAmount += row.Amount
		case "CR":
			daybook.CreditTotalAmount += row.Amount
		}
	}
	daybook.DaybookDetails = daybookDetails
	return daybook, nil
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
			"paymentMethod":   payload.PaymentMethod,
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

func (r daybookRepository) GenerateFinancialStatement(ctx context.Context, company string, year string) ([]mDaybook.DaybookFinancialStatement, error) {
	var accounts []mDaybook.DaybookFinancialStatement
	doc, err := primitive.ObjectIDFromHex(company)
	if err != nil {
		r.logger.Error(err)
	}
	const (
		layoutISO = "2006-01-02T15:04:05.000Z"
	)
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}
	from, _ := time.Parse(layoutISO, fmt.Sprintf("%d-01-01T00:00:00.000Z", yearInt))
	to, _ := time.Parse(layoutISO, fmt.Sprintf("%d-01-01T00:00:00.000Z", yearInt+1))
	filter := bson.M{"company": doc, "daybook.transactionDate": bson.M{
		"$gte": from,
		"$lt":  to,
	}}
	pipeline := []bson.M{
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
		{
			"$lookup": bson.M{
				"from":         "daybooks",
				"localField":   "daybook",
				"foreignField": "_id",
				"as":           "daybook",
			},
		},
		{
			"$match": filter,
		},
		{
			"$unwind": bson.M{
				"path":                       "$daybook",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$group": bson.M{
				"_id": "$account",
				"daybookDetails": bson.M{
					"$addToSet": "$$ROOT",
				},
			},
		},
		{
			"$addFields": bson.M{
				"id":          "$_id._id",
				"code":        "$_id.code",
				"name":        "$_id.name",
				"description": "$_id.description",
				"type":        "$_id.type",
				"company":     "$_id.company",
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
			"$sort": bson.M{"code": 1},
		},
	}

	cur, err := r.Collection.DaybookDetail.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error(err)
	}
	// var test []bson.M
	if err = cur.All(ctx, &accounts); err != nil {
		r.logger.Error(err)
	}
	return accounts, nil
}
