package models

import (
	"log"
	"time"
)

type Event struct {
	DateTime       time.Time     `json:"date_time"`
	Duration       time.Duration `json:"duration"`
	NotifiedBefore time.Duration `json:"notified_before"`
	ID             string        `json:"id"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	UserID         string        `json:"user_id"`
}

func New(generateID func() (string, error)) *Event {
	newID, err := generateID()
	if err != nil {
		log.Fatalln(err)
	}

	return &Event{ID: newID}
}
