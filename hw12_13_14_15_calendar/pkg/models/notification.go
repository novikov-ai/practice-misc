package models

import "time"

type Notification struct {
	EventID string
	Title   string
	UserID  string
	Date    time.Time
}
