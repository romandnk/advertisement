package service

import (
	"github.com/google/uuid"
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
