package models

import "time"

type Image struct {
	ID        string
	Data      []byte
	CreatedAt time.Time
}
