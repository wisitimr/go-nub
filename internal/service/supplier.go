package service

import (
	"context"
	"findigitalservice/internal/auth"
	mRepo "findigitalservice/internal/model/repository"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	mSupplier "findigitalservice/internal/model/supplier"
	mUser "findigitalservice/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type supplierService struct {
	supplierRepo mRepo.SupplierRepository
	logger       *logrus.Logger
}

func InitSupplierService(repo mRepo.Repository, logger *logrus.Logger) mService.SupplierService {
	return &supplierService{
		supplierRepo: repo.Supplier,
		logger:       logger,
	}
}

func (s supplierService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.supplierRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s supplierService) FindAll(ctx context.Context, query map[string][]string) ([]mSupplier.Supplier, error) {
	res, err := s.supplierRepo.FindAll(ctx, query)
	if err != nil {
		return []mSupplier.Supplier{}, err
	}
	return res, nil
}

func (s supplierService) FindById(ctx context.Context, id string) (mSupplier.Supplier, error) {
	res, err := s.supplierRepo.FindById(ctx, id)
	if err != nil {
		return mSupplier.Supplier{}, err
	}
	return res, nil
}

func (s supplierService) Delete(ctx context.Context, id string) error {
	err := s.supplierRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s supplierService) Create(ctx context.Context, payload mSupplier.Supplier) (mSupplier.Supplier, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.supplierRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s supplierService) Update(ctx context.Context, id string, payload mSupplier.Supplier) (mSupplier.Supplier, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mSupplier.Supplier{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.supplierRepo.Update(ctx, payload)
	if err != nil {
		return mSupplier.Supplier{}, err
	}
	return res, nil
}
