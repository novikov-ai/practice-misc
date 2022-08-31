package models

import (
	"log"
	"time"
)

type Event struct {
	DateTime       time.Time
	Duration       time.Duration
	NotifiedBefore time.Duration
	ID             string
	Title          string
	Description    string
	UserID         string
}

func New(generateID func() (string, error)) *Event {
	newID, err := generateID()
	if err != nil {
		log.Fatalln(err)
	}

	return &Event{ID: newID}
}
