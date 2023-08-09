package http

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/romandnk/advertisement/internal/models"
	"net/http"
)

var (
	createUserAction = "create user"
)

type bodyUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userFromBody bodyUser

	err := json.NewDecoder(r.Body).Decode(&userFromBody)
	if err != nil {
		resp := newResponse("", "invalid JSON data", err)
		h.logError(resp.Message, createUserAction, resp.Error)
		renderResponse(w, r, http.StatusBadRequest, resp)
		return
	}

	user := models.User{
		Email:    userFromBody.Email,
		Password: userFromBody.Password,
	}

	id, err := h.service.SignUp(r.Context(), user)
	if err != nil {
		resp := newResponse("", "error creating user", err)
		h.logError(resp.Message, createUserAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"id": id})
}
