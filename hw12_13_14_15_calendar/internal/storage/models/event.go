package models

import (
	"log"
	"time"

	"github.com/google/uuid"
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

func New() *Event {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalln(err)
	}

	return &Event{ID: newUUID.String()}
}
