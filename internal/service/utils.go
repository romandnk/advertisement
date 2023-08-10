package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/romandnk/advertisement/internal/custom_error"
	"github.com/romandnk/advertisement/internal/models"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func comparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createJWT(secret []byte, userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["user_id"] = userID

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func findImageByID(path string, id string) ([]byte, error) {
	var data []byte
	name := id + ".jpg"

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == name {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			dataFile, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			data = dataFile
			return nil
		}

		return nil
	})

	return data, err
}
