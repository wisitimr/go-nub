package service

import (
	"context"
	"saved/http/rest/internal/auth"
	mCompany "saved/http/rest/internal/model/company"
	mRepo "saved/http/rest/internal/model/repository"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type companyService struct {
	companyRepo mRepo.CompanyRepository
	logger      *logrus.Logger
}

func InitCompanyService(companyRepo mRepo.CompanyRepository, logger *logrus.Logger) mService.CompanyService {
	return &companyService{
		companyRepo: companyRepo,
		logger:      logger,
	}
}

func (s companyService) FindAll(ctx context.Context, query map[string][]string) ([]mCompany.Company, error) {
	res, err := s.companyRepo.FindAll(ctx, query)
	if err != nil {
		return []mCompany.Company{}, err
	}
	return res, nil
}

func (s companyService) FindById(ctx context.Context, id string) (mCompany.Company, error) {
	res, err := s.companyRepo.FindById(ctx, id)
	if err != nil {
		return mCompany.Company{}, err
	}
	return res, nil
}

func (s companyService) Create(ctx context.Context, payload mCompany.Company) (mCompany.Company, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.companyRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s companyService) Update(ctx context.Context, id string, payload mCompany.Company) (mCompany.Company, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mCompany.Company{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.companyRepo.Update(ctx, payload)
	if err != nil {
		return mCompany.Company{}, err
	}
	return res, nil
}
