package service

import (
	"context"
	"nub/internal/auth"
	mPaymentMethod "nub/internal/model/payment_method"
	mRepo "nub/internal/model/repository"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentMethodService struct {
	paymentMethodRepo mRepo.PaymentMethodRepository
	logger            *logrus.Logger
}

func InitPaymentMethodService(repo mRepo.Repository, logger *logrus.Logger) mService.PaymentMethodService {
	return &paymentMethodService{
		paymentMethodRepo: repo.PaymentMethod,
		logger:            logger,
	}
}

func (s paymentMethodService) FindAll(ctx context.Context, query map[string][]string) ([]mPaymentMethod.PaymentMethod, error) {
	res, err := s.paymentMethodRepo.FindAll(ctx, query)
	if err != nil {
		return []mPaymentMethod.PaymentMethod{}, err
	}
	return res, nil
}

func (s paymentMethodService) FindById(ctx context.Context, id string) (mPaymentMethod.PaymentMethod, error) {
	res, err := s.paymentMethodRepo.FindById(ctx, id)
	if err != nil {
		return mPaymentMethod.PaymentMethod{}, err
	}
	return res, nil
}

func (s paymentMethodService) Create(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.paymentMethodRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s paymentMethodService) Update(ctx context.Context, id string, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mPaymentMethod.PaymentMethod{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.paymentMethodRepo.Update(ctx, payload)
	if err != nil {
		return mPaymentMethod.PaymentMethod{}, err
	}
	return res, nil
}
