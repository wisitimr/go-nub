package handler

import (
	"encoding/json"
	"net/http"
	mDaybook "saved/http/rest/internal/model/daybook"
	mHandler "saved/http/rest/internal/model/handler"
	mRes "saved/http/rest/internal/model/response"
	mService "saved/http/rest/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type daybookHandler struct {
	daybookService mService.DaybookService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitDaybookHandler(daybookService mService.DaybookService, logger *logrus.Logger) mHandler.DaybookHandler {
	return daybookHandler{
		daybookService: daybookService,
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

func (h daybookHandler) FindAllDetail(w http.ResponseWriter, r *http.Request) {
	daybooks, err := h.daybookService.FindAllDetail(r.Context(), r.URL.Query())
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