package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"path/filepath"
	"time"
)

var pathToImages = "static/images/"

func (h *Handler) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var advert models.Advert

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		resp := newResponse("", "error parsing form", err)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		resp := newResponse("price", "must be a  number e.g. 123.45", err)
		renderResponse(w, r, http.StatusBadRequest, resp)
		return
	}

	imagesForm := r.MultipartForm.File["images"]

	var images []*models.Image
	for _, imageForm := range imagesForm {
		if filepath.Ext(imageForm.Filename) != ".jpg" && filepath.Ext(imageForm.Filename) != ".jpeg" {
			resp := newResponse("images", "image must be with .jpg or .jpeg extension: "+imageForm.Filename, nil)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		var image models.Image

		file, err := imageForm.Open()
		if err != nil {
			file.Close()
			resp := newResponse("images", "error opening file: "+imageForm.Filename, err)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		imageData, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			resp := newResponse("images", "error reading file: "+imageForm.Filename, err)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		image.ID = uuid.New().String()
		image.Data = imageData
		image.CreatedAt = time.Now()

		images = append(images, &image)

		file.Close()
	}

	advert.Title = title
	advert.Description = description
	advert.Price = price
	advert.Images = images

	id, err := h.service.CreateAdvert(r.Context(), advert, pathToImages)
	if err != nil {
		resp := newResponse("", "error creating advert", err)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"id": id})
}

func (h *Handler) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteAdvert(r.Context(), id, pathToImages)
	if err != nil {
		resp := newResponse("", "error deleting advert", err)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusOK)
}
