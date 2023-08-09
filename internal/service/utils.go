package service

import (
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func saveImage(image *models.Image, path string) error {
	err := os.WriteFile(path+image.ID+".jpg", image.Data, 0o644)
	return err
}

func deleteImage(imageID string, path string) error {
	pathImage := path + imageID + ".jpg"
	if _, err := os.Stat(pathImage); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(pathImage)
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 6 {
		return custom_error.CustomError{Field: "password", Message: "min password length is 6"}
	}

	if utf8.RuneCountInString(password) > 18 {
		return custom_error.CustomError{Field: "password", Message: "max password length is 18"}
	}

	if !strings.ContainsAny(password, "0123456789") {
		return custom_error.CustomError{Field: "password", Message: "password must contain at least 1 digit"}
	}

	if strings.ContainsAny(password, " ") {
		return custom_error.CustomError{Field: "password", Message: "password must not contain a space symbol"}
	}

	if !containUpperCaseSymbol(password) {
		return custom_error.CustomError{Field: "password", Message: "password must contain at least 1 upper case symbol"}
	}

	if !containLowerCaseSymbol(password) {
		return custom_error.CustomError{Field: "password", Message: "password must contain at least 1 lower case symbol"}
	}

	return nil
}

func containUpperCaseSymbol(str string) bool {
	for _, i := range str {
		if unicode.IsUpper(i) {
			return true
		}
	}
	return false
}

func containLowerCaseSymbol(str string) bool {
	for _, i := range str {
		if unicode.IsLower(i) {
			return true
		}
	}
	return false
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
