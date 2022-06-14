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

	workersCount := len(tasks)
	if n < workersCount {
		workersCount = n
	}

	for i := 0; i < len(tasks); i += workersCount {
		errs := make(chan error, workersCount)
		wg.Add(workersCount)

		for j := 0; j < workersCount; j++ {
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
