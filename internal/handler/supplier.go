package handler

import (
	"encoding/json"
	"net/http"
	mHandler "nub/internal/model/handler"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"
	mSupplier "nub/internal/model/supplier"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type supplierHandler struct {
	supplierService mService.SupplierService
	logger          *logrus.Logger
	mRes.ResponseDto
}

func InitSupplierHandler(service mService.Service, logger *logrus.Logger) mHandler.SupplierHandler {
	return supplierHandler{
		supplierService: service.Supplier,
		logger:          logger,
	}
}

func (h supplierHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.supplierService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h supplierHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	suppliers, err := h.supplierService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, suppliers, http.StatusOK)
}

func (h supplierHandler) FindById(w http.ResponseWriter, r *http.Request) {
	supplier, err := h.supplierService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, supplier, http.StatusOK)
}

func (h supplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.supplierService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h supplierHandler) Create(w http.ResponseWriter, r *http.Request) {
	supplierPayload := mSupplier.Supplier{}
	err := json.NewDecoder(r.Body).Decode(&supplierPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.supplierService.Create(r.Context(), supplierPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h supplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	supplierPayload := mSupplier.Supplier{}
	err := json.NewDecoder(r.Body).Decode(&supplierPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.supplierService.Update(r.Context(), chi.URLParam(r, "id"), supplierPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
