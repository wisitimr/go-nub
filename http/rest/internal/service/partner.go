package service

import (
	"context"
	"saved/http/rest/internal/auth"
	mPartner "saved/http/rest/internal/model/partner"
	mRepo "saved/http/rest/internal/model/repository"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type partnerService struct {
	partnerRepo mRepo.PartnerRepository
	logger      *logrus.Logger
}

func InitPartnerService(partnerRepo mRepo.PartnerRepository, logger *logrus.Logger) mService.PartnerService {
	return &partnerService{
		partnerRepo: partnerRepo,
		logger:      logger,
	}
}

func (s partnerService) FindAll(ctx context.Context, query map[string][]string) ([]mPartner.Partner, error) {
	res, err := s.partnerRepo.FindAll(ctx, query)
	if err != nil {
		return []mPartner.Partner{}, err
	}
	return res, nil
}

func (s partnerService) FindById(ctx context.Context, id string) (mPartner.Partner, error) {
	res, err := s.partnerRepo.FindById(ctx, id)
	if err != nil {
		return mPartner.Partner{}, err
	}
	return res, nil
}

func (s partnerService) Create(ctx context.Context, payload mPartner.Partner) (mPartner.Partner, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.partnerRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s partnerService) Update(ctx context.Context, id string, payload mPartner.Partner) (mPartner.Partner, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mPartner.Partner{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.partnerRepo.Update(ctx, payload)
	if err != nil {
		return mPartner.Partner{}, err
	}
	return res, nil
}
