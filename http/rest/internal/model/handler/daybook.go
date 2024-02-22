package handler

import (
	"net/http"
)

type DaybookHandler interface {
	Count(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
	GenerateExcel(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}
