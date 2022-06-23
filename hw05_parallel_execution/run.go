package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrErrorsInvalidArgs   = errors.New("invalid arguments")
)

type Task func() error

type ErrorsLimitCalc struct {
	Current, Max int
	Mutex        sync.RWMutex
}

func (errorsCalc *ErrorsLimitCalc) LimitExceeded() bool {
	errorsCalc.Mutex.RLock()
	defer errorsCalc.Mutex.RUnlock()
	return errorsCalc.Current >= errorsCalc.Max
}

func (errorsCalc *ErrorsLimitCalc) IncrementErrors() {
	errorsCalc.Mutex.Lock()
	defer errorsCalc.Mutex.Unlock()

	errorsCalc.Current++
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 || m < 0 {
		return ErrErrorsInvalidArgs
	}

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	jobs := make(chan Task)
	errorsCalc := &ErrorsLimitCalc{Current: 0, Max: m, Mutex: sync.RWMutex{}}
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		go consumer(jobs, errorsCalc, wg)
	}

	defer wg.Wait()
	defer close(jobs)

	for _, task := range tasks {
		jobs <- task
		if errorsCalc.LimitExceeded() {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func consumer(jobs <-chan Task, errorsCalc *ErrorsLimitCalc, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	for job := range jobs {
		err := job()
		if err == nil {
			continue
		}

		errorsCalc.IncrementErrors()
	}
}
