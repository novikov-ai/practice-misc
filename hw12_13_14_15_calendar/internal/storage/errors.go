package storage

import "errors"

var (
	ErrDateBusy                = errors.New("the time is occupied by another event")
	ErrEventAlreadyExists      = errors.New("event with the same UUID already exists")
	ErrEventNotExists          = errors.New("event does not exist")
	ErrInMemoryOperationFailed = errors.New("operation failed")
)
