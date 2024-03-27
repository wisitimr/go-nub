package handler

import (
	"encoding/json"
	mHandler "findigitalservice/internal/model/handler"
	mPartner "findigitalservice/internal/model/partner"
	mRes "findigitalservice/internal/model/response"
	mService "findigitalservice/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type partnerHandler struct {
	partnerService mService.PartnerService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitPartnerHandler(partnerService mService.PartnerService, logger *logrus.Logger) mHandler.PartnerHandler {
	return partnerHandler{
		partnerService: partnerService,
		logger:         logger,
	}
}

func (h partnerHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	partners, err := h.partnerService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, partners, http.StatusOK)
}

func (h partnerHandler) FindById(w http.ResponseWriter, r *http.Request) {
	partner, err := h.partnerService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, partner, http.StatusOK)
}

func (h partnerHandler) Create(w http.ResponseWriter, r *http.Request) {
	partnerPayload := mPartner.Partner{}
	err := json.NewDecoder(r.Body).Decode(&partnerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.partnerService.Create(r.Context(), partnerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h partnerHandler) Update(w http.ResponseWriter, r *http.Request) {
	partnerPayload := mPartner.Partner{}
	err := json.NewDecoder(r.Body).Decode(&partnerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.partnerService.Update(r.Context(), chi.URLParam(r, "id"), partnerPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
