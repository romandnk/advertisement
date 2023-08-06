package http

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/romandnk/advertisement/internal/custom_error"
	"net/http"
)

type response struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func newResponse(field, message string, err error) response {
	var customError custom_error.CustomError

	if errors.As(err, &customError) {
		resp := response{
			Field:   customError.Field,
			Message: message,
			Error:   customError.Message,
		}
		return resp
	}

	return response{
		Field:   field,
		Message: message,
		Error:   err.Error(),
	}
}

func renderResponse(w http.ResponseWriter, r *http.Request, code int, resp response) {
	render.Status(r, code)
	render.JSON(w, r, resp)
}
