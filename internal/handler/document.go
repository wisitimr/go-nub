package handler

import (
	"encoding/json"
	"net/http"
	mDocument "nub/internal/model/document"
	mHandler "nub/internal/model/handler"
	mRes "nub/internal/model/response"
	mService "nub/internal/model/service"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type documentHandler struct {
	documentService mService.DocumentService
	logger          *logrus.Logger
	mRes.ResponseDto
}

func InitDocumentHandler(service mService.Service, logger *logrus.Logger) mHandler.DocumentHandler {
	return documentHandler{
		documentService: service.Document,
		logger:          logger,
	}
}

func (h documentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	documents, err := h.documentService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, documents, http.StatusOK)
}

func (h documentHandler) FindById(w http.ResponseWriter, r *http.Request) {
	document, err := h.documentService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, document, http.StatusOK)
}

func (h documentHandler) Create(w http.ResponseWriter, r *http.Request) {
	documentPayload := mDocument.Document{}
	err := json.NewDecoder(r.Body).Decode(&documentPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.documentService.Create(r.Context(), documentPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h documentHandler) Update(w http.ResponseWriter, r *http.Request) {
	documentPayload := mDocument.Document{}
	err := json.NewDecoder(r.Body).Decode(&documentPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.documentService.Update(r.Context(), chi.URLParam(r, "id"), documentPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
