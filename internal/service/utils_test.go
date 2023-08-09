package service

import (
	"github.com/google/uuid"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestSaveImage(t *testing.T) {
	dir := t.TempDir()

	image := &models.Image{
		ID:        uuid.New().String(),
		Data:      []byte("test"),
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err := saveImage(image, dir)
	require.NoError(t, err)

	_, err = os.Stat(dir + image.ID + ".jpg")
	require.NoError(t, err)
}

func TestDeleteImage(t *testing.T) {
	dir := t.TempDir()

	image := &models.Image{
		ID:        uuid.New().String(),
		Data:      []byte("test"),
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err := saveImage(image, dir)
	require.NoError(t, err)

	err = deleteImage(image.ID, dir)
	require.NoError(t, err)
}

func TestDeleteImageFileIsNotExist(t *testing.T) {
	dir := t.TempDir()

	image := &models.Image{
		ID:        uuid.New().String(),
		Data:      []byte("test"),
		CreatedAt: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	err := saveImage(image, dir)
	require.NoError(t, err)

	err = deleteImage(image.ID, dir)
	require.NoError(t, err)
}

func TestValidatePassword(t *testing.T) {
	testCases := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "Valid password",
			password:      "Abcdef123",
			expectedError: nil,
		},
		{
			name:          "Too short password",
			password:      "Abcd1",
			expectedError: custom_error.CustomError{Field: "password", Message: "min password length is 6"},
		},
		{
			name:          "Too long password",
			password:      "Abcdef123Abcdef1234",
			expectedError: custom_error.CustomError{Field: "password", Message: "max password length is 18"},
		},
		{
			name:          "No digit in password",
			password:      "Abcdefgh",
			expectedError: custom_error.CustomError{Field: "password", Message: "password must contain at least 1 digit"},
		},
		{
			name:          "Password with space",
			password:      "Abc def123",
			expectedError: custom_error.CustomError{Field: "password", Message: "password must not contain a space symbol"},
		},
		{
			name:          "No uppercase symbol in password",
			password:      "abcdef123",
			expectedError: custom_error.CustomError{Field: "password", Message: "password must contain at least 1 upper case symbol"},
		},
		{
			name:          "No lowercase symbol in password",
			password:      "ABCDEF123",
			expectedError: custom_error.CustomError{Field: "password", Message: "password must contain at least 1 lower case symbol"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.password)
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}

func TestContainUpperCaseSymbol(t *testing.T) {
	testCases := []struct {
		name           string
		str            string
		expectedResult bool
	}{
		{
			name:           "Contains uppercase",
			str:            "Abcdef",
			expectedResult: true,
		},
		{
			name:           "Does not contain uppercase",
			str:            "abcdef",
			expectedResult: false,
		},
		{
			name:           "Empty string",
			str:            "",
			expectedResult: false,
		},
		{
			name:           "Only uppercase",
			str:            "ABCDEF",
			expectedResult: true,
		},
		{
			name:           "Mixed case",
			str:            "AbCdEf",
			expectedResult: true,
		},
		{
			name:           "Special characters",
			str:            "AbC!@#$",
			expectedResult: true,
		},
		{
			name:           "Only special characters",
			str:            "!@#$",
			expectedResult: false,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := containUpperCaseSymbol(tc.str)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestContainLowerCaseSymbol(t *testing.T) {
	testCases := []struct {
		name           string
		str            string
		expectedResult bool
	}{
		{
			name:           "Contains lowercase",
			str:            "Abcdef",
			expectedResult: true,
		},
		{
			name:           "Does not contain lowercase",
			str:            "ABCDEF",
			expectedResult: false,
		},
		{
			name:           "Empty string",
			str:            "",
			expectedResult: false,
		},
		{
			name:           "Only lowercase",
			str:            "abcdef",
			expectedResult: true,
		},
		{
			name:           "Mixed case",
			str:            "AbCdEf",
			expectedResult: true,
		},
		{
			name:           "Special characters",
			str:            "AbC!@#$",
			expectedResult: true,
		},
		{
			name:           "Only special characters",
			str:            "!@#$",
			expectedResult: false,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := containLowerCaseSymbol(tc.str)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
