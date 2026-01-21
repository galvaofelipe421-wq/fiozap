package model

import (
	"encoding/json"
	"errors"
	"net/http"

	"fiozap/internal/apperror"
)

type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
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

	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Code: appErr.Code, Error: appErr.Message})
		return
	}
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

func RespondAppError(w http.ResponseWriter, err *apperror.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: err.Code, Error: err.Message})
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
