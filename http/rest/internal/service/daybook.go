package service

import (
	"context"
	"saved/http/rest/internal/auth"
	mDaybook "saved/http/rest/internal/model/daybook"
	mDaybookDetail "saved/http/rest/internal/model/daybook_detail"
	mRepo "saved/http/rest/internal/model/repository"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"
	mUser "saved/http/rest/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type daybookService struct {
	daybookRepo       mRepo.DaybookRepository
	daybookDetailRepo mRepo.DaybookDetailRepository
	logger            *logrus.Logger
}

func InitDaybookService(daybookRepo mRepo.DaybookRepository, daybookDetailRepo mRepo.DaybookDetailRepository, logger *logrus.Logger) mService.DaybookService {
	return &daybookService{
		daybookRepo:       daybookRepo,
		daybookDetailRepo: daybookDetailRepo,
		logger:            logger,
	}
}

func (s daybookService) Count(ctx context.Context, query map[string][]string) (mRes.CountDto, error) {
	count, err := s.daybookRepo.Count(ctx, query)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s daybookService) FindAll(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookList, error) {
	res, err := s.daybookRepo.FindAll(ctx, query)
	if err != nil {
		return []mDaybook.DaybookList{}, err
	}
	return res, nil
}

func (s daybookService) FindAllDetail(ctx context.Context, query map[string][]string) ([]mDaybook.DaybookResponse, error) {
	res, err := s.daybookRepo.FindAllDetail(ctx, query)
	if err != nil {
		return []mDaybook.DaybookResponse{}, err
	}
	return res, nil
}

func (s daybookService) FindById(ctx context.Context, id string) (mDaybook.DaybookResponse, error) {
	res, err := s.daybookRepo.FindById(ctx, id)
	if err != nil {
		return mDaybook.DaybookResponse{}, err
	}
	return res, nil
}

func (s daybookService) Create(ctx context.Context, payload mDaybook.DaybookPayload) (mDaybook.DaybookPayload, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	// new daybook
	newId := primitive.NewObjectID()
	newDt := time.Now()
	inv := mDaybook.Daybook{}
	inv.Id = newId
	inv.Number = payload.Number
	inv.Invoice = payload.Invoice
	inv.Document = payload.Document
	inv.TransactionDate = payload.TransactionDate
	inv.Company = payload.Company
	inv.Supplier = payload.Supplier
	inv.Customer = payload.Customer
	inv.CreatedBy = user.Id
	inv.CreatedAt = newDt
	inv.UpdatedBy = user.Id
	inv.UpdatedAt = newDt

	// new daybook detail
	payload.Id = inv.Id
	payload.CreatedBy = inv.CreatedBy
	payload.CreatedAt = inv.CreatedAt
	payload.UpdatedBy = inv.UpdatedBy
	payload.UpdatedAt = inv.UpdatedAt
	if len(payload.DaybookDetails) > 0 {
		docs := make([]interface{}, len(payload.DaybookDetails))
		for i, doc := range payload.DaybookDetails {
			newId = primitive.NewObjectID()
			newDt = time.Now()
			docs[i] = mDaybookDetail.DaybookDetail{
				Id:        newId,
				Name:      doc.Name,
				Type:      doc.Type,
				Amount:    doc.Amount,
				Account:   doc.Account,
				Daybook:   payload.Id,
				CreatedBy: user.Id,
				CreatedAt: newDt,
				UpdatedBy: user.Id,
				UpdatedAt: newDt,
			}
			inv.DaybookDetails = append(inv.DaybookDetails, newId)
			payload.DaybookDetails[i].Id = newId
			payload.DaybookDetails[i].CreatedBy = user.Id
			payload.DaybookDetails[i].CreatedAt = newDt
			payload.DaybookDetails[i].UpdatedBy = user.Id
			payload.DaybookDetails[i].UpdatedAt = newDt
		}
		err = s.daybookRepo.BulkCreateDaybookDetail(ctx, docs)
		if err != nil {
			return payload, err
		}
	}

	_, err = s.daybookRepo.Create(ctx, inv)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func (s daybookService) Update(ctx context.Context, id string, payload mDaybook.Daybook) (mDaybook.Daybook, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDaybook.Daybook{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.daybookRepo.Update(ctx, payload)
	if err != nil {
		return mDaybook.Daybook{}, err
	}
	return res, nil
}
