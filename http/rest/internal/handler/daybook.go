package handler

import (
	"encoding/json"
	mDaybook "findigitalservice/http/rest/internal/model/daybook"
	mHandler "findigitalservice/http/rest/internal/model/handler"
	mRes "findigitalservice/http/rest/internal/model/response"
	mService "findigitalservice/http/rest/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type daybookHandler struct {
	daybookService mService.DaybookService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitDaybookHandler(service mService.Service, logger *logrus.Logger) mHandler.DaybookHandler {
	return daybookHandler{
		daybookService: service.Daybook,
		logger:         logger,
	}
}

func (h daybookHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.daybookService.Count(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h daybookHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	daybooks, err := h.daybookService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, daybooks, http.StatusOK)
}

func (h daybookHandler) FindById(w http.ResponseWriter, r *http.Request) {
	daybook, err := h.daybookService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, daybook, http.StatusOK)
}

func (h daybookHandler) GenerateExcel(w http.ResponseWriter, r *http.Request) {
	f, err := h.daybookService.GenerateExcel(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}

	h.Respond(w, r, f, 0)
}

func (h daybookHandler) Create(w http.ResponseWriter, r *http.Request) {
	daybookPayload := mDaybook.DaybookPayload{}
	err := json.NewDecoder(r.Body).Decode(&daybookPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.daybookService.Create(r.Context(), daybookPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h daybookHandler) Update(w http.ResponseWriter, r *http.Request) {
	daybookPayload := mDaybook.Daybook{}
	err := json.NewDecoder(r.Body).Decode(&daybookPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.daybookService.Update(r.Context(), chi.URLParam(r, "id"), daybookPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}

func (h daybookHandler) GenerateFinancialStatement(w http.ResponseWriter, r *http.Request) {
	f, err := h.daybookService.GenerateFinancialStatement(r.Context(), chi.URLParam(r, "company"), chi.URLParam(r, "year"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, f, http.StatusOK)
}
