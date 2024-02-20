package service

import (
	"context"
	"saved/http/rest/internal/auth"
	mMaterial "saved/http/rest/internal/model/material"
	mRepo "saved/http/rest/internal/model/repository"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type materialService struct {
	materialRepo mRepo.MaterialRepository
	logger       *logrus.Logger
}

func InitMaterialService(materialRepo mRepo.MaterialRepository, logger *logrus.Logger) mService.MaterialService {
	return &materialService{
		materialRepo: materialRepo,
		logger:       logger,
	}
}

func (s materialService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.materialRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s materialService) FindAll(ctx context.Context, query map[string][]string) ([]mMaterial.Material, error) {
	res, err := s.materialRepo.FindAll(ctx, query)
	if err != nil {
		return []mMaterial.Material{}, err
	}
	return res, nil
}

func (s materialService) FindById(ctx context.Context, id string) (mMaterial.Material, error) {
	res, err := s.materialRepo.FindById(ctx, id)
	if err != nil {
		return mMaterial.Material{}, err
	}
	return res, nil
}

func (s materialService) Delete(ctx context.Context, id string) error {
	err := s.materialRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s materialService) Create(ctx context.Context, payload mMaterial.Material) (mMaterial.Material, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.materialRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s materialService) Update(ctx context.Context, id string, payload mMaterial.Material) (mMaterial.Material, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mMaterial.Material{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.materialRepo.Update(ctx, payload)
	if err != nil {
		return mMaterial.Material{}, err
	}
	return res, nil
}
