package handler

import (
	"encoding/json"
	mCustomer "findigitalservice/internal/model/customer"
	mHandler "findigitalservice/internal/model/handler"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type customerHandler struct {
	customerService mService.CustomerService
	logger          *logrus.Logger
	mRes.ResponseDto
}

func InitCustomerHandler(service mService.Service, logger *logrus.Logger) mHandler.CustomerHandler {
	return customerHandler{
		customerService: service.Customer,
		logger:          logger,
	}
}

func (h customerHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.customerService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h customerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	customers, err := h.customerService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, customers, http.StatusOK)
}

func (h customerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	customer, err := h.customerService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, customer, http.StatusOK)
}

func (h customerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.customerService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h customerHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerPayload := mCustomer.Customer{}
	err := json.NewDecoder(r.Body).Decode(&customerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.customerService.Create(r.Context(), customerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h customerHandler) Update(w http.ResponseWriter, r *http.Request) {
	customerPayload := mCustomer.Customer{}
	err := json.NewDecoder(r.Body).Decode(&customerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.customerService.Update(r.Context(), chi.URLParam(r, "id"), customerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
