package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var getImageAction = "get image by id"

func (h *Handler) GetImageByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	image, err := h.service.GetImageByID(r.Context(), id)
	if err != nil {
		resp := newResponse("", "error getting image by id", err)
		h.logError(resp.Message, getImageAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	w.Header().Set("content-type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(len(image.Data)))
	if _, err := w.Write(image.Data); err != nil {
		resp := newResponse("", "error displaying image", err)
		h.logError(resp.Message, getImageAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}
}
