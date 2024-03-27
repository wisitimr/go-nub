package handler

import (
	"encoding/json"
	mHandler "findigitalservice/internal/model/handler"
	mRes "findigitalservice/internal/model/response"
	mRole "findigitalservice/internal/model/role"
	mService "findigitalservice/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type roleHandler struct {
	roleService mService.RoleService
	logger      *logrus.Logger
	mRes.ResponseDto
}

func InitRoleHandler(service mService.Service, logger *logrus.Logger) mHandler.RoleHandler {
	return roleHandler{
		roleService: service.Role,
		logger:      logger,
	}
}

func (h roleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	roles, err := h.roleService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, roles, http.StatusOK)
}

func (h roleHandler) FindById(w http.ResponseWriter, r *http.Request) {
	role, err := h.roleService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, role, http.StatusOK)
}

func (h roleHandler) Create(w http.ResponseWriter, r *http.Request) {
	rolePayload := mRole.Role{}
	err := json.NewDecoder(r.Body).Decode(&rolePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.roleService.Create(r.Context(), rolePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h roleHandler) Update(w http.ResponseWriter, r *http.Request) {
	rolePayload := mRole.Role{}
	err := json.NewDecoder(r.Body).Decode(&rolePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.roleService.Update(r.Context(), chi.URLParam(r, "id"), rolePayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
