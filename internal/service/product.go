package service

import (
	"context"
	"nub/internal/auth"
	mProduct "nub/internal/model/product"
	mRepo "nub/internal/model/repository"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productService struct {
	productRepo mRepo.ProductRepository
	logger      *logrus.Logger
}

func InitProductService(repo mRepo.Repository, logger *logrus.Logger) mService.ProductService {
	return &productService{
		productRepo: repo.Product,
		logger:      logger,
	}
}

func (s productService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.productRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s productService) FindAll(ctx context.Context, query map[string][]string) ([]mProduct.Product, error) {
	res, err := s.productRepo.FindAll(ctx, query)
	if err != nil {
		return []mProduct.Product{}, err
	}
	return res, nil
}

func (s productService) FindById(ctx context.Context, id string) (mProduct.Product, error) {
	res, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return mProduct.Product{}, err
	}
	return res, nil
}

func (s productService) Delete(ctx context.Context, id string) error {
	err := s.productRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s productService) Create(ctx context.Context, payload mProduct.Product) (mProduct.Product, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.productRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s productService) Update(ctx context.Context, id string, payload mProduct.Product) (mProduct.Product, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mProduct.Product{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.productRepo.Update(ctx, payload)
	if err != nil {
		return mProduct.Product{}, err
	}
	return res, nil
}
