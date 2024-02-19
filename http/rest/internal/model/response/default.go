package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type ResponseDto struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"statusMessage"`
	Data       interface{} `json:"data"`
}

func (s ResponseDto) Respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	var res ResponseDto
	switch v := data.(type) {
	case error:
		res.StatusCode = http.StatusInternalServerError
		if v.Error() == mongo.ErrNoDocuments.Error() {
			res.StatusCode = http.StatusNotFound
		} else if v.Error() == "incorrect email or password, please try again. " {
			res.StatusCode = http.StatusBadRequest
		} else if strings.HasPrefix(v.Error(), "this email address") {
			res.StatusCode = http.StatusBadRequest
		}
		res.Message = v.Error()
	default:
		res.StatusCode = http.StatusOK
		res.Message = "Success"
		res.Data = data
	}
	if status != 0 {
		res.StatusCode = status
	}
	method := strings.ToLower(r.Method)
	if (method == "post" || method == "put") && (res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		res.Message = "ทำรายการสำเร็จ"
	}
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}
