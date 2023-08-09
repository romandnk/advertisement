package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	mock_logger "github.com/romandnk/advertisement/internal/logger/mock"
	"github.com/romandnk/advertisement/internal/models"
	mock_service "github.com/romandnk/advertisement/internal/service/mock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

const urlAdverts = "/api/v1/adverts"

func TestHandlerCreateAdvert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockServices(ctrl)

	expectedID := uuid.New().String()
	expectedUserID := uuid.New().String()

	expectedAdvert := models.Advert{
		Title:       "test create",
		Description: "test create",
		Price:       decimal.New(1200, 0),
		UserID:      expectedUserID,
	}

	ctx := context.Background()

	services.EXPECT().CreateAdvert(gomock.Any(), expectedAdvert).Return(expectedID, nil)

	handler := NewHandler(services, nil)

	r := chi.NewRouter()
	r.Post(urlAdverts, handler.CreateAdvert)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	err := bodyWriter.WriteField("title", "test create")
	require.NoError(t, err)
	err = bodyWriter.WriteField("description", "test create")
	require.NoError(t, err)
	err = bodyWriter.WriteField("price", "1200")
	require.NoError(t, err)
	file, err := bodyWriter.CreateFormFile("image", "test.jpg")
	require.NoError(t, err)
	n, err := file.Write([]byte("test image"))
	require.NoError(t, err)
	require.Equal(t, len([]byte("test image")), n)

	err = bodyWriter.Close()
	require.NoError(t, err)

	w := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAdverts, bodyBuf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	ctx = context.WithValue(ctx, "user_id", expectedUserID)

	r.ServeHTTP(w, req.WithContext(ctx))

	require.Equal(t, http.StatusCreated, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	id, ok := responseBody["id"]
	require.Equal(t, expectedID, id)
	require.True(t, ok)
}

func TestHandlerCreateAdvertError(t *testing.T) {
	testCases := []struct {
		name           string
		message        string
		expectedError  string
		expectedAdvert map[string]string
		code           int
		contentType    string
		responseBody   map[string]interface{}
	}{
		{
			name:          "content type is application/json",
			message:       "error parsing form",
			expectedError: "request Content-Type isn't multipart/form-data",
			code:          http.StatusInternalServerError,
			contentType:   "application/json",
			responseBody: map[string]interface{}{
				"message": "error parsing form",
				"error":   "request Content-Type isn't multipart/form-data",
			},
		},
		{
			name:          "invalid price",
			message:       "must be a  number e.g. 123.45",
			expectedError: "can't convert invalid price to decimal: exponent is not numeric",
			expectedAdvert: map[string]string{
				"price": "invalid price",
			},
			code:        http.StatusBadRequest,
			contentType: "multipart/form-data",
			responseBody: map[string]interface{}{
				"field":   "price",
				"message": "must be a  number e.g. 123.45",
				"error":   "can't convert invalid price to decimal: exponent is not numeric",
			},
		},
		{
			name:          "empty user id",
			message:       "invalid user id ctx",
			expectedError: "",
			expectedAdvert: map[string]string{
				"price": "1200",
			},
			code:        http.StatusInternalServerError,
			contentType: "multipart/form-data",
			responseBody: map[string]interface{}{
				"field":   "user_id",
				"message": "invalid user id ctx",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(tc.message,
				zap.String("action", createAdvertAction),
				zap.String("error", tc.expectedError),
			)

			handler := NewHandler(services, logger)
			r := chi.NewRouter()
			r.Post(urlAdverts, handler.CreateAdvert)

			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			err := bodyWriter.WriteField("title", tc.expectedAdvert["title"])
			require.NoError(t, err)
			err = bodyWriter.WriteField("description", tc.expectedAdvert["description"])
			require.NoError(t, err)
			err = bodyWriter.WriteField("price", tc.expectedAdvert["price"])
			require.NoError(t, err)
			file, err := bodyWriter.CreateFormFile("image", "test.jpg")
			require.NoError(t, err)
			n, err := file.Write([]byte("test image"))
			require.NoError(t, err)
			require.Equal(t, len([]byte("test image")), n)

			ctx := context.Background()

			err = bodyWriter.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAdverts, bodyBuf)
			require.NoError(t, err)
			if tc.name != "empty user id" {
				ctx = context.WithValue(ctx, "user_id", "test_user")
			}

			switch tc.contentType {
			case "multipart/form-data":
				req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
			case "application/json":
				req.Header.Set("Content-Type", tc.contentType)
			}

			r.ServeHTTP(w, req.WithContext(ctx))

			require.Equal(t, tc.code, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			require.Equal(t, tc.responseBody, responseBody)
		})
	}
}

func TestHandlerCreateAdvertErrorCreatingAdvert(t *testing.T) {
	message := "error creating advert"
	testCases := []struct {
		name           string
		expectedAdvert models.Advert
		expectedError  error
		responseBody   map[string]interface{}
	}{
		{
			name: "title is empty",
			expectedAdvert: models.Advert{
				Title:       "",
				Description: "test",
				Price:       decimal.New(1200, 0),
				UserID:      "test_user",
			},
			expectedError: custom_error.CustomError{
				Field:   "title",
				Message: "empty title",
			},
			responseBody: map[string]interface{}{
				"field":   "title",
				"message": message,
				"error":   "empty title",
			},
		},
		{
			name: "negative price",
			expectedAdvert: models.Advert{
				Title:       "test",
				Description: "test",
				Price:       decimal.New(-1200, 0),
				UserID:      "test_user",
			},
			expectedError: custom_error.CustomError{
				Field:   "price",
				Message: "negative price",
			},
			responseBody: map[string]interface{}{
				"field":   "price",
				"message": message,
				"error":   "negative price",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			services.EXPECT().CreateAdvert(gomock.Any(), tc.expectedAdvert).Return("", tc.expectedError)
			logger.EXPECT().Error(message,
				zap.String("action", createAdvertAction),
				zap.String("error", tc.expectedError.Error()),
			)

			handler := NewHandler(services, logger)
			r := chi.NewRouter()
			r.Post(urlAdverts, handler.CreateAdvert)

			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			err := bodyWriter.WriteField("title", tc.expectedAdvert.Title)
			require.NoError(t, err)
			err = bodyWriter.WriteField("description", tc.expectedAdvert.Description)
			require.NoError(t, err)
			err = bodyWriter.WriteField("price", tc.expectedAdvert.Price.String())
			require.NoError(t, err)
			file, err := bodyWriter.CreateFormFile("image", "test.jpg")
			require.NoError(t, err)
			n, err := file.Write([]byte("test"))
			require.NoError(t, err)
			require.Equal(t, len([]byte("test")), n)

			ctx := context.Background()

			err = bodyWriter.Close()
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAdverts, bodyBuf)
			require.NoError(t, err)
			ctx = context.WithValue(ctx, "user_id", "test_user")

			req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

			r.ServeHTTP(w, req.WithContext(ctx))

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			require.Equal(t, tc.responseBody, responseBody)
		})
	}
}

func TestHandlerDeleteAdvert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	services := mock_service.NewMockServices(ctrl)

	expectedID := uuid.New().String()

	ctx := context.Background()

	services.EXPECT().DeleteAdvert(gomock.Any(), expectedID).Return(nil)

	handler := NewHandler(services, nil)

	r := chi.NewRouter()
	r.Post(urlAdverts+"/{id}", handler.DeleteAdvert)

	w := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAdverts+"/"+expectedID, nil)
	require.NoError(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandlerDeleteAdvertError(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		expectedError error
		responseBody  map[string]interface{}
	}{
		{
			name:          "invalid id length",
			id:            "test id",
			expectedError: custom_error.CustomError{Field: "id", Message: "invalid UUID length: 7"},
			responseBody: map[string]interface{}{
				"field":   "id",
				"message": "error deleting advert",
				"error":   "invalid UUID length: 7",
			},
		},
		{
			name:          "invalid id format",
			id:            "1be0349d-cc15-452a-9b8d999-b1d7bd-e0",
			expectedError: custom_error.CustomError{Field: "id", Message: "invalid UUID format"},
			responseBody: map[string]interface{}{
				"field":   "id",
				"message": "error deleting advert",
				"error":   "invalid UUID format",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			ctx := context.Background()

			services.EXPECT().DeleteAdvert(gomock.Any(), tc.id).Return(tc.expectedError)
			logger.EXPECT().Error("error deleting advert",
				zap.String("action", deleteAdvertAction),
				zap.String("error", tc.expectedError.Error()),
			)

			handler := NewHandler(services, logger)

			r := chi.NewRouter()
			r.Post(urlAdverts+"/{id}", handler.DeleteAdvert)

			w := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAdverts+"/"+tc.id, nil)
			require.NoError(t, err)

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			require.Equal(t, tc.responseBody, responseBody)
		})
	}
}
