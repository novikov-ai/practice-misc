package models

type User struct {
	ID     string
	Events map[string]*Event
}
