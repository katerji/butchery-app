package httpresponse

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error any    `json:"error,omitempty"`
	Meta  any    `json:"meta,omitempty"`
}

func JSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Response{Data: data})
}

func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Response{Data: data})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{Error: message})
}

func ValidationError(w http.ResponseWriter, details any) {
	JSON(w, http.StatusUnprocessableEntity, Response{Error: details})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
