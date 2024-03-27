package handler

import (
	"encoding/json"
	mCompany "findigitalservice/internal/model/company"
	mHandler "findigitalservice/internal/model/handler"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type companyHandler struct {
	companyService mService.CompanyService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitCompanyHandler(service mService.Service, logger *logrus.Logger) mHandler.CompanyHandler {
	return companyHandler{
		companyService: service.Company,
		logger:         logger,
	}
}

func (h companyHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	companys, err := h.companyService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, companys, http.StatusOK)
}

func (h companyHandler) FindById(w http.ResponseWriter, r *http.Request) {
	company, err := h.companyService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, company, http.StatusOK)
}

func (h companyHandler) Create(w http.ResponseWriter, r *http.Request) {
	companyPayload := mCompany.Company{}
	err := json.NewDecoder(r.Body).Decode(&companyPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.companyService.Create(r.Context(), companyPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h companyHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyPayload := mCompany.Company{}
	err := json.NewDecoder(r.Body).Decode(&companyPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.companyService.Update(r.Context(), chi.URLParam(r, "id"), companyPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
