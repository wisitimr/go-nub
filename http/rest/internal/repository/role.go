package repository

import (
	"context"
	mCollection "saved/http/rest/internal/model/collection"
	mRepo "saved/http/rest/internal/model/repository"
	mRole "saved/http/rest/internal/model/role"
	"saved/http/rest/internal/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roleRepository struct {
	Collection mCollection.Collection
	logger     *logrus.Logger
}

func InitRoleRepository(collection mCollection.Collection, logger *logrus.Logger) mRepo.RoleRepository {
	return &roleRepository{
		Collection: collection,
		logger:     logger,
	}
}

func (r roleRepository) FindAll(ctx context.Context, query map[string][]string) ([]mRole.Role, error) {
	roles := []mRole.Role{}
	cur, err := r.Collection.Role.Find(ctx, util.QueryHandler(query))
	if err != nil {
		return roles, err
	}
	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var e mRole.Role
		err := cur.Decode(&e)
		if err != nil {
			r.logger.Fatal(err)
		}
		roles = append(roles, e)
	}
	return roles, nil
}

func (r roleRepository) FindById(ctx context.Context, id string) (mRole.Role, error) {
	role := mRole.Role{}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return role, err
	}
	err = r.Collection.Role.FindOne(ctx, bson.M{"_id": doc}).Decode(&role)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r roleRepository) Create(ctx context.Context, payload mRole.Role) (mRole.Role, error) {
	if _, err := r.Collection.Role.InsertOne(ctx, payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r roleRepository) Update(ctx context.Context, payload mRole.Role) (mRole.Role, error) {
	var updated mRole.Role
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
	err := r.Collection.Role.FindOneAndUpdate(ctx, filter, update, &opt).Decode(&updated)
	if err != nil {
		return updated, err
	}
	return updated, nil
}
