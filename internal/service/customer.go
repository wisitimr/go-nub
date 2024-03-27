package service

import (
	"context"
	"findigitalservice/internal/auth"
	mCustomer "findigitalservice/internal/model/customer"
	mRepo "findigitalservice/internal/model/repository"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	mUser "findigitalservice/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customerService struct {
	customerRepo mRepo.CustomerRepository
	logger       *logrus.Logger
}

func InitCustomerService(repo mRepo.Repository, logger *logrus.Logger) mService.CustomerService {
	return &customerService{
		customerRepo: repo.Customer,
		logger:       logger,
	}
}

func (s customerService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.customerRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s customerService) FindAll(ctx context.Context, query map[string][]string) ([]mCustomer.Customer, error) {
	res, err := s.customerRepo.FindAll(ctx, query)
	if err != nil {
		return []mCustomer.Customer{}, err
	}
	return res, nil
}

func (s customerService) FindById(ctx context.Context, id string) (mCustomer.Customer, error) {
	res, err := s.customerRepo.FindById(ctx, id)
	if err != nil {
		return mCustomer.Customer{}, err
	}
	return res, nil
}

func (s customerService) Delete(ctx context.Context, id string) error {
	err := s.customerRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s customerService) Create(ctx context.Context, payload mCustomer.Customer) (mCustomer.Customer, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.customerRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s customerService) Update(ctx context.Context, id string, payload mCustomer.Customer) (mCustomer.Customer, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mCustomer.Customer{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.customerRepo.Update(ctx, payload)
	if err != nil {
		return mCustomer.Customer{}, err
	}
	return res, nil
}
