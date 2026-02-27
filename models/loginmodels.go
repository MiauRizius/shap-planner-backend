package models

import "time"

type RefreshToken struct {
	ID        string    `json:id`
	UserID    string    `json:userid`
	Token     string    `json:token`
	ExpiresAt time.Time `json:expiresat`
}
