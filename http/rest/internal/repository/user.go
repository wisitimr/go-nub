package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mRepo "saved/http/rest/internal/model/repository"
	mUser "saved/http/rest/internal/model/user"
	"saved/http/rest/internal/util"
	"sort"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitUserRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.UserRepository {
	return &userRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r userRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.Collection.User.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error(err)
	}
	return count, nil
}

func (r userRepository) FindAll(ctx context.Context, query map[string][]string) ([]mUser.User, error) {
	users := []mUser.User{}
	cur, err := r.Collection.User.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return users, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mUser.User
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		users = append(users, e)
	}
	return users, nil
}

// func (r userRepository) FindById(ctx context.Context, id string) (mUser.UserCompany, error) {
// 	users := []mUser.UserCompany{}
// 	doc, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return mUser.UserCompany{}, err
// 	}
// 	pipeline := []bson.M{
// 		{
// 			"$match": bson.M{"_id": doc},
// 		},
// 		{
// 			"$lookup": bson.M{
// 				// Define the details collection for the join.
// 				"from": "companies",
// 				// Specify the variable to use in the pipeline stage.
// 				"let": bson.M{
// 					"companies": "$companies",
// 				},
// 				"pipeline": []bson.M{
// 					// Select only the relevant details from the details collection.
// 					// Otherwise all the details are selected.
// 					{
// 						"$match": bson.M{
// 							"$expr": bson.M{
// 								"$in": []interface{}{
// 									"$_id",
// 									"$$companies",
// 								},
// 							},
// 						},
// 					},
// 					// Sort details by their createdAt field in asc. -1 = desc
// 					{
// 						"$sort": bson.M{
// 							"createdAt": 1,
// 						},
// 					},
// 				},
// 				// Use details as the field name to match struct field.
// 				"as": "companies",
// 			},
// 		},
// 	}

// 	cur, err := r.Collection.User.Aggregate(ctx, pipeline)
// 	if err != nil {
// 		r.logger.Error(err)
// 	}
// 	if err = cur.All(ctx, &users); err != nil {
// 		r.logger.Error(err)
// 	}
// 	return users[0], nil
// }

func (r userRepository) FindById(ctx context.Context, id string) (mUser.UserCompany, error) {
	user := mUser.User{}
	result := mUser.UserCompany{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mUser.UserCompany{}, err
	}
	err = r.Collection.User.FindOne(ctx, bson.M{"_id": doc}).Decode(&user)
	if err != nil {
		return result, err
	}
	result.Id = user.Id
	result.Username = user.Username
	result.FirstName = user.FirstName
	result.LastName = user.LastName
	result.Email = user.Email
	result.Role = user.Role
	result.CreatedBy = user.CreatedBy
	result.CreatedAt = user.CreatedAt
	result.UpdatedBy = user.UpdatedBy
	result.UpdatedAt = user.UpdatedAt
	outCompanies := []mUser.Company{}
	ch := make(chan mUser.Company)
	for _, doc := range user.Companies {
		go func(ch chan mUser.Company, doc primitive.ObjectID) {
			var out mUser.Company
			// err = r.DaybookDetail.FindOne(ctx, bson.M{"_id": doc}).Decode(&out)
			err = r.Collection.Company.FindOne(ctx, bson.M{"_id": doc}).Decode(&out)
			ch <- out
		}(ch, doc)
	}
	for range user.Companies {
		outCompanies = append(outCompanies, <-ch)
	}
	sort.Slice(outCompanies[:], func(i, j int) bool {
		return outCompanies[i].CreatedAt.Before(outCompanies[j].CreatedAt)
	})
	result.Companies = outCompanies
	return result, nil
}

func (r userRepository) FindUserProfile(ctx context.Context, doc primitive.ObjectID) (mUser.UserProfile, error) {
	user := mUser.User{}
	result := mUser.UserProfile{}
	err := r.Collection.User.FindOne(ctx, bson.M{"_id": doc}).Decode(&user)
	if err != nil {
		return result, err
	}
	result.Id = user.Id
	result.Username = user.Username
	result.FirstName = user.FirstName
	result.LastName = user.LastName
	result.Email = user.Email
	result.Role = user.Role
	outCompanies := []mUser.Company{}
	ch := make(chan mUser.Company)
	for _, doc := range user.Companies {
		go func(ch chan mUser.Company, doc primitive.ObjectID) {
			var out mUser.Company
			// err = r.DaybookDetail.FindOne(ctx, bson.M{"_id": doc}).Decode(&out)
			err = r.Collection.Company.FindOne(ctx, bson.M{"_id": doc}).Decode(&out)
			ch <- out
		}(ch, doc)
	}
	for range user.Companies {
		outCompanies = append(outCompanies, <-ch)
	}
	sort.Slice(outCompanies[:], func(i, j int) bool {
		return outCompanies[i].CreatedAt.Before(outCompanies[j].CreatedAt)
	})
	result.Companies = outCompanies
	return result, nil
}

func (r userRepository) FindUserCompany(ctx context.Context, doc primitive.ObjectID) (mUser.UserCompany, error) {
	users := []mUser.UserCompany{}
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": doc},
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
				// Define the details collection for the join.
				"from": "companies",
				// Specify the variable to use in the pipeline stage.
				"let": bson.M{
					"companies": "$companies",
				},
				"pipeline": []bson.M{
					// Select only the relevant details from the details collection.
					// Otherwise all the details are selected.
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$in": []interface{}{
									"$_id",
									"$$companies",
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
				},
				// Use details as the field name to match struct field.
				"as": "companies",
			},
		},
	}

	cur, err := r.Collection.User.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Error(err)
	}
	if err = cur.All(ctx, &users); err != nil {
		r.logger.Error(err)
	}
	return users[0], nil
}

func (r userRepository) FindByUsername(ctx context.Context, username string) (mUser.User, error) {
	user := mUser.User{}
	err := r.Collection.User.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r userRepository) FindByEmail(ctx context.Context, email string) (mUser.User, error) {
	user := mUser.User{}
	err := r.Collection.User.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r userRepository) Create(ctx context.Context, payload mUser.User) (mUser.User, error) {
	if _, err := r.Collection.User.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r userRepository) Update(ctx context.Context, payload mUser.User) (mUser.User, error) {
	var updated mUser.User
	filter := bson.M{"_id": payload.Id}
	update := bson.M{
		"$set": bson.M{
			"firstName": payload.FirstName,
			"lastName":  payload.LastName,
			"email":     payload.Email,
			"companies": payload.Companies,
			// "role":      payload.Role,
			"updatedBy": payload.UpdatedBy,
			"updatedAt": payload.UpdatedAt,
		},
	}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := r.Collection.User.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
