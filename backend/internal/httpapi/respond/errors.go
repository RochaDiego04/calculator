package respond

import "net/http"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type ErrorBody struct {
	Errors []Error `json:"errors"`
}

func Errors(w http.ResponseWriter, status int, errs ...Error) {
	JSON(w, status, ErrorBody{Errors: errs})
}
