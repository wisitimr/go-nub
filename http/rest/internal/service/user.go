package service

import (
	"context"
	"errors"
	"findigitalservice/http/rest/internal/auth"
	mRepo "findigitalservice/http/rest/internal/model/repository"
	mRes "findigitalservice/http/rest/internal/model/response"
	mService "findigitalservice/http/rest/internal/model/service"
	mUser "findigitalservice/http/rest/internal/model/user"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo mRepo.UserRepository
	logger   *logrus.Logger
}

func InitUserService(userRepo mRepo.UserRepository, logger *logrus.Logger) mService.UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s userService) Count(ctx context.Context) (mRes.CountDto, error) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return mRes.CountDto{Count: 0}, err
	}
	return mRes.CountDto{Count: count}, nil
}

func (s userService) FindAll(ctx context.Context, query map[string][]string) ([]mUser.User, error) {
	res, err := s.userRepo.FindAll(ctx, query)
	if err != nil {
		return []mUser.User{}, err
	}
	return res, nil
}

func (s userService) FindById(ctx context.Context, id string) (mUser.UserCompany, error) {
	res, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return mUser.UserCompany{}, err
	}
	return res, nil
}

func (s userService) FindUserProfile(ctx context.Context) (mUser.UserProfile, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	res, err := s.userRepo.FindUserProfile(ctx, user.Id)
	if err != nil {
		return mUser.UserProfile{}, err
	}
	res.FullName = res.FirstName + " " + res.LastName
	return res, nil
}

func (s userService) FindUserCompany(ctx context.Context) (mUser.UserCompany, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	res, err := s.userRepo.FindUserCompany(ctx, user.Id)
	if err != nil {
		return mUser.UserCompany{}, err
	}
	return res, nil
}

func (s userService) Create(ctx context.Context, payload mUser.User) (mUser.User, error) {
	_, err := s.userRepo.FindByUsername(ctx, payload.Username)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		newId := primitive.NewObjectID()
		user := mUser.User{
			Id:        newId,
			Username:  payload.Username,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  auth.Hash([]byte(payload.Password)),
			CreatedBy: newId,
			CreatedAt: time.Now(),
			UpdatedBy: newId,
			UpdatedAt: time.Now(),
		}
		res, err := s.userRepo.Create(ctx, user)
		if err != nil {
			return res, err
		}
		return res, nil
	} else {
		return mUser.User{}, fmt.Errorf("this username \"%s\" is already in use", payload.Username)
	}
}

func (s userService) Update(ctx context.Context, id string, payload mUser.User) (mUser.UpdatedUserProfile, error) {
	user, err := auth.UserLogin(ctx, s.logger)
	if err != nil {
		user = mUser.User{}
	}
	doc, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return mUser.UpdatedUserProfile{}, err
	}
	payload.Id = doc
	payload.UpdatedBy = user.Id
	payload.UpdatedAt = time.Now()
	res, err := s.userRepo.Update(ctx, payload)
	if err != nil {
		return mUser.UpdatedUserProfile{}, err
	}
	u := mUser.UpdatedUserProfile{
		Id:        res.Id,
		Username:  res.Username,
		FullName:  res.FirstName + " " + res.LastName,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Companies: res.Companies,
		Role:      res.Role,
	}
	return u, nil
}

func (s userService) Login(ctx context.Context, payload mUser.Login) (mRes.TokenDto, error) {
	user, err := s.userRepo.FindByUsername(ctx, payload.Username)
	if err != nil {
		s.logger.Error(err)
		return mRes.TokenDto{}, errors.New("incorrect username or password, please try again. ")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return mRes.TokenDto{}, errors.New("incorrect username or password, please try again. ")
	}
	jwtToken, err := auth.GenerateToken(user)
	if err != nil {
		return mRes.TokenDto{}, err
	}
	s.logger.Info("jwtToken : ", jwtToken)
	return jwtToken, nil
}
