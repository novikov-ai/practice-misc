package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrErrorsInvalidArgs = errors.New("arguments invalid")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	if n*len(tasks) == 0 {
		return ErrErrorsInvalidArgs
	}

	jobs := make(chan Task)
	go producer(tasks, jobs)

	done := make(chan struct{})
	errors := make(chan error)

	for i := 0; i < n; i++ {
		go consumer(done, jobs, errors)
	}

	completed := 0
	for {
		select {
		case <-errors:
			m--
			if m == 0 {
				return ErrErrorsLimitExceeded
			}

		case <-done:
			completed++
			if completed == n {
				return nil
			}
		}
	}

	return nil
}

func producer(tasks []Task, jobs chan<- Task) {
	for _, task := range tasks {
		jobs <- task
	}
	close(jobs)
}

func consumer(done chan struct{}, jobs <-chan Task, errors chan error) {
	for job := range jobs {
		err := job()
		if err != nil {
			errors <- err
		}
	}
	done <- struct{}{}
}
