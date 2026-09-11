package httpapi

import "net/http"

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type errorBody struct {
	Errors []apiError `json:"errors"`
}

func writeErrors(w http.ResponseWriter, status int, errs ...apiError) {
	writeJSON(w, status, errorBody{Errors: errs})
}
