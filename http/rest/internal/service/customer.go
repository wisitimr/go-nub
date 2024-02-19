package service

import (
	"context"
	"saved/http/rest/internal/auth"
	mCustomer "saved/http/rest/internal/model/customer"
	mRepo "saved/http/rest/internal/model/repository"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customerService struct {
	customerRepo mRepo.CustomerRepository
	logger       *logrus.Logger
}

func InitCustomerService(customerRepo mRepo.CustomerRepository, logger *logrus.Logger) mService.CustomerService {
	return &customerService{
		customerRepo: customerRepo,
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
