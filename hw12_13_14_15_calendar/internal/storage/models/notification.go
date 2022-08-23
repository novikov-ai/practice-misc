package models

import "time"

type Notification struct {
	Date   time.Time
	ID     string
	Title  string
	UserID string
}
