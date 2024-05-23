package domain

import "time"

type Promo struct {
	ID          string
	UserID      string
	PromoTypeID string
	PromoCode   string
	Status      int
	ValidUntil  time.Time
	CreatedAt   time.Time
}
