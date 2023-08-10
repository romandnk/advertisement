package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"time"
)

var (
	createAdvertAction  = "create advert"
	deleteAdvertAction  = "delete advert"
	getAdvertByIDAction = "get advert by id"
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
			file.Close()
			resp := newResponse("images", "image cannot be decoded: "+imageForm.Filename, nil)
			h.logError(resp.Message, createAdvertAction, resp.Error)
			renderResponse(w, r, http.StatusBadRequest, resp)
			return
		}

		file.Seek(0, 0)

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
	userID := r.Context().Value("user_id")
	switch userID.(type) {
	case string:
		advert.UserID = userID.(string)
	default:
		resp := newResponse("user_id", "invalid user id ctx", err)
		h.logError(resp.Message, createAdvertAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}
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

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetAdvertByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	advert, err := h.service.GetAdvertByID(r.Context(), id)
	if err != nil {
		resp := newResponse("", "error getting advert by id", err)
		h.logError(resp.Message, getAdvertByIDAction, resp.Error)
		renderResponse(w, r, http.StatusInternalServerError, resp)
		return
	}

	var imageURLs []string
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	for _, img := range advert.Images {
		url := fmt.Sprintf("http://%s:%s/api/v1/images/%s", host, port, img.ID)
		imageURLs = append(imageURLs, url)
	}

	jsonResponse := struct {
		ID          string          `json:"id""`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Price       decimal.Decimal `json:"price"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
		UserID      string          `json:"user_id"`
		ImageURLs   []string        `json:"image_urls"`
	}{
		ID:          advert.ID,
		Title:       advert.Title,
		Description: advert.Description,
		Price:       advert.Price,
		CreatedAt:   advert.CreatedAt,
		UpdatedAt:   advert.UpdatedAt,
		UserID:      advert.UserID,
		ImageURLs:   imageURLs,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, jsonResponse)
}
