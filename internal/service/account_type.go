package service

import (
	"context"
	"nub/internal/auth"
	mAccountType "nub/internal/model/account_type"
	mRepo "nub/internal/model/repository"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type accountTypeService struct {
	accountTypeRepo mRepo.AccountTypeRepository
	logger          *logrus.Logger
}

func InitAccountTypeService(repo mRepo.Repository, logger *logrus.Logger) mService.AccountTypeService {
	return &accountTypeService{
		accountTypeRepo: repo.AccountType,
		logger:          logger,
	}
}

func (s accountTypeService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.accountTypeRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s accountTypeService) FindAll(ctx context.Context, query map[string][]string) ([]mAccountType.AccountType, error) {
	res, err := s.accountTypeRepo.FindAll(ctx, query)
	if err != nil {
		return []mAccountType.AccountType{}, err
	}
	return res, nil
}

func (s accountTypeService) FindById(ctx context.Context, id string) (mAccountType.AccountType, error) {
	res, err := s.accountTypeRepo.FindById(ctx, id)
	if err != nil {
		return mAccountType.AccountType{}, err
	}
	return res, nil
}

func (s accountTypeService) Delete(ctx context.Context, id string) error {
	err := s.accountTypeRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s accountTypeService) Create(ctx context.Context, payload mAccountType.AccountType) (mAccountType.AccountType, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.accountTypeRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s accountTypeService) Update(ctx context.Context, id string, payload mAccountType.AccountType) (mAccountType.AccountType, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mAccountType.AccountType{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.accountTypeRepo.Update(ctx, payload)
	if err != nil {
		return mAccountType.AccountType{}, err
	}
	return res, nil
}
