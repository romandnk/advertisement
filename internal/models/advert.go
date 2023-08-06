package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Advert struct {
	ID          string
	Title       string
	Description string
	Price       decimal.Decimal
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      string
	Deleted     bool
	Images      []*Image
}
