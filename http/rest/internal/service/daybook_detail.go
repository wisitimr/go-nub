package service

import (
	"context"
	"findigitalservice/http/rest/internal/auth"
	mDaybook "findigitalservice/http/rest/internal/model/daybook"
	mDaybookDetail "findigitalservice/http/rest/internal/model/daybook_detail"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	mService "findigitalservice/http/rest/internal/model/service"
	mUser "findigitalservice/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type daybookDetailService struct {
	daybookRepo       mRepo.DaybookRepository
	daybookDetailRepo mRepo.DaybookDetailRepository
	logger            *logrus.Logger
}

func InitDaybookDetailService(repo mRepo.Repository, logger *logrus.Logger) mService.DaybookDetailService {
	return &daybookDetailService{
		daybookRepo:       repo.Daybook,
		daybookDetailRepo: repo.DaybookDetail,
		logger:            logger,
	}
}

func (s daybookDetailService) FindAll(ctx context.Context, query map[string][]string) ([]mDaybookDetail.DaybookDetail, error) {
	res, err := s.daybookDetailRepo.FindAll(ctx, query)
	if err != nil {
		return []mDaybookDetail.DaybookDetail{}, err
	}
	return res, nil
}

func (s daybookDetailService) FindById(ctx context.Context, id string) (mDaybookDetail.DaybookDetail, error) {
	res, err := s.daybookDetailRepo.FindById(ctx, id)
	if err != nil {
		return mDaybookDetail.DaybookDetail{}, err
	}
	return res, nil
}

func (s daybookDetailService) Create(ctx context.Context, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.daybookDetailRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	daybook, err := s.daybookRepo.FindById(ctx, res.Daybook.Hex())
	if err != nil {
		return res, err
	}
	s.logger.Error(daybook.DaybookDetails)
	var daybookDetails []primitive.ObjectID
	if daybook.DaybookDetails != nil {
		for _, detail := range daybook.DaybookDetails {
			daybookDetails = append(daybookDetails, detail.Id)
		}
	}
	daybookDetails = append(daybookDetails, res.Id)
	var newDaybook mDaybook.Daybook
	newDaybook.Id = daybook.Id
	newDaybook.Number = daybook.Number
	newDaybook.Invoice = daybook.Invoice
	newDaybook.Document = daybook.Document
	newDaybook.TransactionDate = daybook.TransactionDate
	newDaybook.Company = daybook.Company
	newDaybook.Supplier = daybook.Supplier
	newDaybook.Customer = daybook.Customer
	newDaybook.PaymentMethod = daybook.PaymentMethod
	newDaybook.DaybookDetails = daybookDetails
	newDaybook.UpdatedBy = user.Id
	newDaybook.UpdatedAt = time.Now()
	s.daybookRepo.Update(ctx, newDaybook)
	return res, nil
}

func (s daybookDetailService) Update(ctx context.Context, id string, payload mDaybookDetail.DaybookDetail) (mDaybookDetail.DaybookDetail, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDaybookDetail.DaybookDetail{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.daybookDetailRepo.Update(ctx, payload)
	if err != nil {
		return mDaybookDetail.DaybookDetail{}, err
	}
	return res, nil
}
