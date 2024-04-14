package handler

import (
	"encoding/json"
	"net/http"
	mHandler "nub/internal/model/handler"
	mPaymentMethod "nub/internal/model/payment_method"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type paymentMethodHandler struct {
	paymentMethodService mService.PaymentMethodService
	logger               *logrus.Logger
	mRes.ResponseDto
}

func InitPaymentMethodHandler(service mService.Service, logger *logrus.Logger) mHandler.PaymentMethodHandler {
	return paymentMethodHandler{
		paymentMethodService: service.PaymentMethod,
		logger:               logger,
	}
}

func (h paymentMethodHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	paymentMethods, err := h.paymentMethodService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, paymentMethods, http.StatusOK)
}

func (h paymentMethodHandler) FindById(w http.ResponseWriter, r *http.Request) {
	paymentMethod, err := h.paymentMethodService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, paymentMethod, http.StatusOK)
}

func (h paymentMethodHandler) Create(w http.ResponseWriter, r *http.Request) {
	paymentMethodPayload := mPaymentMethod.PaymentMethod{}
	err := json.NewDecoder(r.Body).Decode(&paymentMethodPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.paymentMethodService.Create(r.Context(), paymentMethodPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h paymentMethodHandler) Update(w http.ResponseWriter, r *http.Request) {
	paymentMethodPayload := mPaymentMethod.PaymentMethod{}
	err := json.NewDecoder(r.Body).Decode(&paymentMethodPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.paymentMethodService.Update(r.Context(), chi.URLParam(r, "id"), paymentMethodPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
