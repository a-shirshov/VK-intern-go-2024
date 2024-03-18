package response

import (
	"encoding/json"
	"net/http"
	"vk-intern_test-case/models"
)

func MakeJsonEncoder(w http.ResponseWriter) *json.Encoder {
	w.Header().Set("Content-Type", "application/json")
	jsonEnc := json.NewEncoder(w)
	return jsonEnc
}

func WriteResponse(w http.ResponseWriter, jsonEnc *json.Encoder, httpStatus int, response any) {
	w.WriteHeader(httpStatus)
	jsonEnc.Encode(response)
}

func WriteBasicResponse(w http.ResponseWriter, jsonEnc *json.Encoder, httpStatus int, message string) {
	w.WriteHeader(httpStatus)
	jsonEnc.Encode(&models.BasicResponse{Status: message})
}
