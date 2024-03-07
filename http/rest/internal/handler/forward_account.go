package handler

import (
	"encoding/json"
	mForwardAccount "findigitalservice/http/rest/internal/model/forward_account"
	mHandler "findigitalservice/http/rest/internal/model/handler"
	mRes "findigitalservice/http/rest/internal/model/response"
	mService "findigitalservice/http/rest/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type forwardAccountHandler struct {
	forwardAccountService mService.ForwardAccountService
	logger                *logrus.Logger
	mRes.ResponseDto
}

func InitForwardAccountHandler(service mService.Service, logger *logrus.Logger) mHandler.ForwardAccountHandler {
	return forwardAccountHandler{
		forwardAccountService: service.ForwardAccount,
		logger:                logger,
	}
}

func (h forwardAccountHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.forwardAccountService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h forwardAccountHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	forwardAccounts, err := h.forwardAccountService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, forwardAccounts, http.StatusOK)
}

func (h forwardAccountHandler) FindById(w http.ResponseWriter, r *http.Request) {
	forwardAccount, err := h.forwardAccountService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, forwardAccount, http.StatusOK)
}

func (h forwardAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.forwardAccountService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h forwardAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	forwardAccountPayload := mForwardAccount.ForwardAccount{}
	err := json.NewDecoder(r.Body).Decode(&forwardAccountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.forwardAccountService.Create(r.Context(), forwardAccountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h forwardAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	forwardAccountPayload := mForwardAccount.ForwardAccount{}
	err := json.NewDecoder(r.Body).Decode(&forwardAccountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.forwardAccountService.Update(r.Context(), chi.URLParam(r, "id"), forwardAccountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
