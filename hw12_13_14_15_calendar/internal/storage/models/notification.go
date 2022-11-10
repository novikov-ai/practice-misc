package models

import "time"

type Notification struct {
	Date   time.Time `json:"date"`
	ID     string    `json:"id"`
	Title  string    `json:"title"`
	UserID string    `json:"user_id"`
}
