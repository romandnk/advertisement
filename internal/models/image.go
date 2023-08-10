package models

import "time"

type Image struct {
	ID        string
	Data      []byte
	AdvertID  string
	CreatedAt time.Time
	Deleted   bool
}
