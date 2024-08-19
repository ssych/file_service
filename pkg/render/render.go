package render

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	fallbackMarshalErr = `{"error": "failed to marshal error message"}`
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Success(
	w http.ResponseWriter,
	body interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func Error(
	w http.ResponseWriter,
	statusCode int,
	err error,
) {
	w.Header().Set("Content-Type", "application/json")

	response := &ErrorResponse{
		Error: err.Error(),
	}

	js, merr := json.Marshal(response)
	if merr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		_, _ = io.WriteString(w, fallbackMarshalErr)
		return
	}

	w.WriteHeader(statusCode)
	if _, merr = w.Write(js); merr != nil {
		http.Error(w, merr.Error(), http.StatusInternalServerError)
		return
	}
}
