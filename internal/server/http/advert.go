package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/shopspring/decimal"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
)

var (
	createAdvertAction = "create advert"
	deleteAdvertAction = "delete advert"
)

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var advert models.Advert

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		resp := newResponse("", "error parsing form", err)
		h.logError(resp.Message, createAdvertAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		resp := newResponse("price", "must be a  number e.g. 123.45", err)
		h.logError(resp.Message, createAdvertAction, resp.Error)
		renderResponse(w, r, http.StatusBadRequest, resp)
		return
	}

	imagesForm := r.MultipartForm.File["images"]

	var images []*models.Image
	for _, imageForm := range imagesForm {
		var img models.Image

		file, err := imageForm.Open()
		if err != nil {
			file.Close()
			resp := newResponse("images", "error opening file: "+imageForm.Filename, err)
			h.logError(resp.Message, createAdvertAction, resp.Error)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		_, _, err = image.Decode(file)
		if err != nil {
			resp := newResponse("images", "image cannot be decoded: "+imageForm.Filename, nil)
			h.logError(resp.Message, createAdvertAction, resp.Error)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		imageData, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			resp := newResponse("images", "error reading file: "+imageForm.Filename, err)
			h.logError(resp.Message, createAdvertAction, resp.Error)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		img.Data = imageData

		images = append(images, &img)

		file.Close()
	}

	advert.Title = title
	advert.Description = description
	advert.Price = price
	advert.Images = images

	id, err := h.service.CreateAdvert(r.Context(), advert)
	if err != nil {
		resp := newResponse("", "error creating advert", err)
		h.logError(resp.Message, createAdvertAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"id": id})
}

func (h *Handler) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteAdvert(r.Context(), id)
	if err != nil {
		resp := newResponse("", "error deleting advert", err)
		h.logError(resp.Message, deleteAdvertAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusOK)
}
