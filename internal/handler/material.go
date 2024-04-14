package handler

import (
	"encoding/json"
	"net/http"
	mHandler "nub/internal/model/handler"
	mMaterial "nub/internal/model/material"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type materialHandler struct {
	materialService mService.MaterialService
	logger          *logrus.Logger
	mRes.ResponseDto
}

func InitMaterialHandler(service mService.Service, logger *logrus.Logger) mHandler.MaterialHandler {
	return materialHandler{
		materialService: service.Material,
		logger:          logger,
	}
}

func (h materialHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.materialService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h materialHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	materials, err := h.materialService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, materials, http.StatusOK)
}

func (h materialHandler) FindById(w http.ResponseWriter, r *http.Request) {
	material, err := h.materialService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, material, http.StatusOK)
}

func (h materialHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.materialService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h materialHandler) Create(w http.ResponseWriter, r *http.Request) {
	materialPayload := mMaterial.Material{}
	err := json.NewDecoder(r.Body).Decode(&materialPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.materialService.Create(r.Context(), materialPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h materialHandler) Update(w http.ResponseWriter, r *http.Request) {
	materialPayload := mMaterial.Material{}
	err := json.NewDecoder(r.Body).Decode(&materialPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.materialService.Update(r.Context(), chi.URLParam(r, "id"), materialPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
