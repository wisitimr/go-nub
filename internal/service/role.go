package service

import (
	"context"
	"findigitalservice/internal/auth"
	mRepo "findigitalservice/internal/model/repository"
	mRole "findigitalservice/internal/model/role"
	mService "findigitalservice/internal/model/service"
	mUser "findigitalservice/internal/model/user"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roleService struct {
	roleRepo mRepo.RoleRepository
	logger   *logrus.Logger
}

func InitRoleService(repo mRepo.Repository, logger *logrus.Logger) mService.RoleService {
	return &roleService{
		roleRepo: repo.Role,
		logger:   logger,
	}
}

func (s roleService) FindAll(ctx context.Context, query map[string][]string) ([]mRole.Role, error) {
	res, err := s.roleRepo.FindAll(ctx, query)
	if err != nil {
		return []mRole.Role{}, err
	}
	return res, nil
}

func (s roleService) FindById(ctx context.Context, id string) (mRole.Role, error) {
	res, err := s.roleRepo.FindById(ctx, id)
	if err != nil {
		return mRole.Role{}, err
	}
	return res, nil
}

func (s roleService) Create(ctx context.Context, payload mRole.Role) (mRole.Role, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	payload.Id = primitive.NewObjectID()
	payload.CreatedBy = user.Id
	payload.CreatedAt = time.Now()
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.roleRepo.Create(ctx, payload)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s roleService) Update(ctx context.Context, id string, payload mRole.Role) (mRole.Role, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mRole.Role{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.roleRepo.Update(ctx, payload)
	if err != nil {
		return mRole.Role{}, err
	}
	return res, nil
}
