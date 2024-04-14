package service

import (
	"context"
	"nub/internal/auth"
	mDocument "nub/internal/model/document"
	mRepo "nub/internal/model/repository"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type documentService struct {
	documentRepo mRepo.DocumentRepository
	logger       *logrus.Logger
}

func InitDocumentService(repo mRepo.Repository, logger *logrus.Logger) mService.DocumentService {
	return &documentService{
		documentRepo: repo.Document,
		logger:       logger,
	}
}

func (s documentService) FindAll(ctx context.Context, query map[string][]string) ([]mDocument.Document, error) {
	res, err := s.documentRepo.FindAll(ctx, query)
	if err != nil {
		return []mDocument.Document{}, err
	}
	return res, nil
}

func (s documentService) FindById(ctx context.Context, id string) (mDocument.Document, error) {
	res, err := s.documentRepo.FindById(ctx, id)
	if err != nil {
		return mDocument.Document{}, err
	}
	return res, nil
}

func (s documentService) Create(ctx context.Context, payload mDocument.Document) (mDocument.Document, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.documentRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s documentService) Update(ctx context.Context, id string, payload mDocument.Document) (mDocument.Document, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mDocument.Document{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.documentRepo.Update(ctx, payload)
	if err != nil {
		return mDocument.Document{}, err
	}
	return res, nil
}
