package handler

import (
	"encoding/json"
	"net/http"
	mAccountType "nub/internal/model/account_type"
	mHandler "nub/internal/model/handler"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type accountTypeHandler struct {
	accountTypeService mService.AccountTypeService
	logger             *logrus.Logger
	mRes.ResponseDto
}

func InitAccountTypeHandler(service mService.Service, logger *logrus.Logger) mHandler.AccountTypeHandler {
	return accountTypeHandler{
		accountTypeService: service.AccountType,
		logger:             logger,
	}
}

func (h accountTypeHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.accountTypeService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h accountTypeHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	accountTypes, err := h.accountTypeService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, accountTypes, http.StatusOK)
}

func (h accountTypeHandler) FindById(w http.ResponseWriter, r *http.Request) {
	accountType, err := h.accountTypeService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, accountType, http.StatusOK)
}

func (h accountTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.accountTypeService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h accountTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountTypePayload := mAccountType.AccountType{}
	err := json.NewDecoder(r.Body).Decode(&accountTypePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.accountTypeService.Create(r.Context(), accountTypePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h accountTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountTypePayload := mAccountType.AccountType{}
	err := json.NewDecoder(r.Body).Decode(&accountTypePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.accountTypeService.Update(r.Context(), chi.URLParam(r, "id"), accountTypePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
