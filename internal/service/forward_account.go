package service

import (
	"context"
	"nub/internal/auth"
	mForwardAccount "nub/internal/model/forward_account"
	mRepo "nub/internal/model/repository"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type forwardAccountService struct {
	forwardAccountRepo mRepo.ForwardAccountRepository
	logger             *logrus.Logger
}

func InitForwardAccountService(repo mRepo.Repository, logger *logrus.Logger) mService.ForwardAccountService {
	return &forwardAccountService{
		forwardAccountRepo: repo.ForwardAccount,
		logger:             logger,
	}
}

func (s forwardAccountService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.forwardAccountRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s forwardAccountService) FindAll(ctx context.Context, query map[string][]string) ([]mForwardAccount.ForwardAccount, error) {
	res, err := s.forwardAccountRepo.FindAll(ctx, query)
	if err != nil {
		return []mForwardAccount.ForwardAccount{}, err
	}
	return res, nil
}

func (s forwardAccountService) FindById(ctx context.Context, id string) (mForwardAccount.ForwardAccount, error) {
	res, err := s.forwardAccountRepo.FindById(ctx, id)
	if err != nil {
		return mForwardAccount.ForwardAccount{}, err
	}
	return res, nil
}

func (s forwardAccountService) Delete(ctx context.Context, id string) error {
	err := s.forwardAccountRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s forwardAccountService) Create(ctx context.Context, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.forwardAccountRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s forwardAccountService) Update(ctx context.Context, id string, payload mForwardAccount.ForwardAccount) (mForwardAccount.ForwardAccount, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mForwardAccount.ForwardAccount{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.forwardAccountRepo.Update(ctx, payload)
	if err != nil {
		return mForwardAccount.ForwardAccount{}, err
	}
	return res, nil
}
