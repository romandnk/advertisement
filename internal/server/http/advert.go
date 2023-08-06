package http

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

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

	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")

	var images []*models.Image
	for _, imageForm := range imagesForm {
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

		url := fmt.Sprintf("http://%s:%d/api/v1/images/%s", host, port, image.ID)
		image.Url = url

		images = append(images, &image)

		file.Close()
	}

	advert.Title = title
	advert.Description = description
	advert.Price = price
	advert.Images = images

	id, err := h.service.CreateAdvert(r.Context(), advert)
	if err != nil {
		resp := newResponse("", "error creating advert", err)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{"id": id})
}
