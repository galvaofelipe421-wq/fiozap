package model

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

func RespondOK(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusOK, data)
}

func RespondCreated(w http.ResponseWriter, data interface{}) {
	RespondJSON(w, http.StatusCreated, data)
}

func RespondBadRequest(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusBadRequest, err)
}

func RespondUnauthorized(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusUnauthorized, err)
}

func RespondForbidden(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusForbidden, err)
}

func RespondNotFound(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusNotFound, err)
}

func RespondInternalError(w http.ResponseWriter, err error) {
	RespondError(w, http.StatusInternalServerError, err)
}
