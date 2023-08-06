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
	Images      []*Image
}
