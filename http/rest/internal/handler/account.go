package handler

import (
	"encoding/json"
	"net/http"
	mAccount "saved/http/rest/internal/model/account"
	mHandler "saved/http/rest/internal/model/handler"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type accountHandler struct {
	accountService mService.AccountService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitAccountHandler(accountService mService.AccountService, logger *logrus.Logger) mHandler.AccountHandler {
	return accountHandler{
		accountService: accountService,
		logger:         logger,
	}
}

func (h accountHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.accountService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h accountHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, accounts, http.StatusOK)
}

func (h accountHandler) FindById(w http.ResponseWriter, r *http.Request) {
	account, err := h.accountService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, account, http.StatusOK)
}

func (h accountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.accountService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h accountHandler) Create(w http.ResponseWriter, r *http.Request) {
	accountPayload := mAccount.Account{}
	err := json.NewDecoder(r.Body).Decode(&accountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.accountService.Create(r.Context(), accountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h accountHandler) Update(w http.ResponseWriter, r *http.Request) {
	accountPayload := mAccount.Account{}
	err := json.NewDecoder(r.Body).Decode(&accountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.accountService.Update(r.Context(), chi.URLParam(r, "id"), accountPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
