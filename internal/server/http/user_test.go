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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

const urlUsers = "/api/v1/users"

func TestHandlerSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock_service.NewMockServices(ctrl)

	expectedUser := models.User{
		Email:    "test@vk.com",
		Password: "Qwerty123",
	}

	expectedID := uuid.New().String()

	service.EXPECT().SignUp(gomock.Any(), expectedUser).Return(expectedID, nil)

	handler := NewHandler(service, nil, secretTest)

	r := chi.NewRouter()
	r.Post(urlUsers+"/sign-up", handler.SignUp)

	expectedBody := map[string]string{
		"email":    expectedUser.Email,
		"password": expectedUser.Password,
	}

	jsonBody, err := json.Marshal(&expectedBody)
	require.NoError(t, err)

	ctx := context.Background()

	w := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlUsers+"/sign-up", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	id, ok := responseBody["id"]
	require.Equal(t, expectedID, id)
	require.True(t, ok)
}

func TestHandlerSignUpError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_logger.NewMockLogger(ctrl)

	logger.EXPECT().Error("invalid JSON data",
		zap.String("action", createUserAction),
		zap.String("error", "json: cannot unmarshal bool into Go struct field bodyUser.email of type string"),
	)

	handler := NewHandler(nil, logger, secretTest)

	r := chi.NewRouter()
	r.Post(urlUsers+"/sign-up", handler.SignUp)

	expectedBody := map[string]interface{}{
		"email":    true,
		"password": "password",
	}

	jsonBody, err := json.Marshal(&expectedBody)
	require.NoError(t, err)

	ctx := context.Background()

	w := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlUsers+"/sign-up", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	expectedResponse := map[string]interface{}{
		"message": "invalid JSON data",
		"error":   "json: cannot unmarshal bool into Go struct field bodyUser.email of type string",
	}

	require.Equal(t, expectedResponse, responseBody)
}

func TestHandlerSignUpErrorCreatingUser(t *testing.T) {
	testCases := []struct {
		name             string
		expectedUser     models.User
		expectedError    custom_error.CustomError
		expectedBody     map[string]interface{}
		expectedResponse map[string]interface{}
	}{
		{
			name: "empty email",
			expectedUser: models.User{
				Email:    "",
				Password: "Qwerty123",
			},
			expectedError: custom_error.CustomError{
				Field:   "email",
				Message: "mail: no address",
			},
			expectedBody: map[string]interface{}{
				"email":    "",
				"password": "Qwerty123",
			},
			expectedResponse: map[string]interface{}{
				"field":   "email",
				"message": "error creating user",
				"error":   "mail: no address",
			},
		},

		{
			name: "empty password",
			expectedUser: models.User{
				Email:    "test@vk.com",
				Password: "",
			},
			expectedError: custom_error.CustomError{
				Field:   "password",
				Message: "min password length is 6",
			},
			expectedBody: map[string]interface{}{
				"email":    "test@vk.com",
				"password": "",
			},
			expectedResponse: map[string]interface{}{
				"field":   "password",
				"message": "error creating user",
				"error":   "min password length is 6",
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			service.EXPECT().SignUp(gomock.Any(), tc.expectedUser).Return("", tc.expectedError)
			logger.EXPECT().Error("error creating user",
				zap.String("action", createUserAction),
				zap.String("error", tc.expectedError.Error()),
			)

			handler := NewHandler(service, logger, secretTest)

			r := chi.NewRouter()
			r.Post(urlUsers+"/sign-up", handler.SignUp)

			jsonBody, err := json.Marshal(&tc.expectedBody)

			require.NoError(t, err)

			ctx := context.Background()

			w := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlUsers+"/sign-up", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusInternalServerError, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, responseBody)
		})
	}
}
