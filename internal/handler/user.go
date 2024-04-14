package handler

import (
	"encoding/json"
	"net/http"
	mHandler "nub/internal/model/handler"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"
	mUser "nub/internal/model/user"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type userHandler struct {
	userService mService.UserService
	logger      *logrus.Logger
	mRes.ResponseDto
}

func InitUserHandler(service mService.Service, logger *logrus.Logger) mHandler.UserHandler {
	return userHandler{
		userService: service.User,
		logger:      logger,
	}
}

func (h userHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.userService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h userHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, users, http.StatusOK)
}

func (h userHandler) FindById(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, user, http.StatusOK)
}

func (h userHandler) FindUserProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.FindUserProfile(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, user, http.StatusOK)
}

func (h userHandler) FindUserCompany(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.FindUserCompany(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, user, http.StatusOK)
}

func (h userHandler) Create(w http.ResponseWriter, r *http.Request) {
	userPayload := mUser.User{}
	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.userService.Create(r.Context(), userPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h userHandler) Update(w http.ResponseWriter, r *http.Request) {
	userPayload := mUser.User{}
	err := json.NewDecoder(r.Body).Decode(&userPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.userService.Update(r.Context(), chi.URLParam(r, "id"), userPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}

func (h userHandler) Login(w http.ResponseWriter, r *http.Request) {
	payload := mUser.Login{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.userService.Login(r.Context(), payload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
