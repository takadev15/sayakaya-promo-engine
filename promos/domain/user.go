package domain

import "time"

type User struct {
	ID          string
	Name        string
	Email       string
	PhoneNumber string
	Status      bool
	Birthday    time.Time
}
