package domain

import "time"

type PromoType struct {
	ID        string
	Name      string
	Rule      string
	CreatedBy string
	CreatedAt time.Time
}
