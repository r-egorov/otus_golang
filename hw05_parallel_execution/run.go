package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errsCount := 0

	for i := 0; i < len(tasks); i += n {
		errs := make(chan error, n)
		wg.Add(n)

		for j := 0; j < n; j++ {
			go func(task_idx int) {
				defer wg.Done()
				errs <- tasks[task_idx]()
			}(i + j)
		}

		wg.Wait()
		close(errs)

		for err := range errs {
			if err != nil {
				errsCount++
			}
		}

		if errsCount >= m {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
