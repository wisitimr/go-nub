package handler

import (
	"encoding/json"
	mDaybookDetail "findigitalservice/internal/model/daybook_detail"
	mHandler "findigitalservice/internal/model/handler"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type DaybookDetailHandler struct {
	DaybookDetailService mService.DaybookDetailService
	logger               *logrus.Logger
	mRes.ResponseDto
}

func InitDaybookDetailHandler(service mService.Service, logger *logrus.Logger) mHandler.DaybookDetailHandler {
	return DaybookDetailHandler{
		DaybookDetailService: service.DaybookDetail,
		logger:               logger,
	}
}

func (h DaybookDetailHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	DaybookDetails, err := h.DaybookDetailService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, DaybookDetails, http.StatusOK)
}

func (h DaybookDetailHandler) FindById(w http.ResponseWriter, r *http.Request) {
	DaybookDetail, err := h.DaybookDetailService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, DaybookDetail, http.StatusOK)
}

func (h DaybookDetailHandler) Create(w http.ResponseWriter, r *http.Request) {
	DaybookDetailPayload := mDaybookDetail.DaybookDetail{}
	err := json.NewDecoder(r.Body).Decode(&DaybookDetailPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.DaybookDetailService.Create(r.Context(), DaybookDetailPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h DaybookDetailHandler) Update(w http.ResponseWriter, r *http.Request) {
	DaybookDetailPayload := mDaybookDetail.DaybookDetail{}
	err := json.NewDecoder(r.Body).Decode(&DaybookDetailPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.DaybookDetailService.Update(r.Context(), chi.URLParam(r, "id"), DaybookDetailPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
