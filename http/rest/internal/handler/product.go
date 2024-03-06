package handler

import (
	"encoding/json"
	mHandler "findigitalservice/http/rest/internal/model/handler"
	mProduct "findigitalservice/http/rest/internal/model/product"
	mRes "findigitalservice/http/rest/internal/model/response"
	mService "findigitalservice/http/rest/internal/model/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type productHandler struct {
	productService mService.ProductService
	logger         *logrus.Logger
	mRes.ResponseDto
}

func InitProductHandler(service mService.Service, logger *logrus.Logger) mHandler.ProductHandler {
	return productHandler{
		productService: service.Product,
		logger:         logger,
	}
}

func (h productHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.productService.Count(r.Context())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, count, http.StatusOK)
}

func (h productHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.FindAll(r.Context(), r.URL.Query())
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, products, http.StatusOK)
}

func (h productHandler) FindById(w http.ResponseWriter, r *http.Request) {
	product, err := h.productService.FindById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, product, http.StatusOK)
}

func (h productHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.productService.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, nil, http.StatusOK)
}

func (h productHandler) Create(w http.ResponseWriter, r *http.Request) {
	productPayload := mProduct.Product{}
	err := json.NewDecoder(r.Body).Decode(&productPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.productService.Create(r.Context(), productPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusCreated)
}

func (h productHandler) Update(w http.ResponseWriter, r *http.Request) {
	productPayload := mProduct.Product{}
	err := json.NewDecoder(r.Body).Decode(&productPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	res, err := h.productService.Update(r.Context(), chi.URLParam(r, "id"), productPayload)
	if err != nil {
		h.Respond(w, r, err, 0)
		return
	}
	h.Respond(w, r, res, http.StatusOK)
}
