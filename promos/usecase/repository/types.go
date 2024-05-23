package repository

import "time"

type BirthdayPromo struct {
	PromoCode   string    `db:"promo_code"`
	UserID      string    `db:"user_id"`
	PromoTypeID string    `db:"promo_type_id"`
	ValidUntil  time.Time `db:"valid_until"`
	Status      int       `db:"status"`
}

type PromoType struct {
	PromoTypeID string `db:"id"`
	Name        string `db:"name"`
	Rule        string `db:"rule"`
}
