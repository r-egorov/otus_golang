package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// If m <= 0, it ignores all the errors.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errsCount := 0

	workersCount := len(tasks) // launch no more goroutines than tasks
	if n < workersCount {
		workersCount = n
	}

	errs := make(chan error, workersCount)
	defer close(errs)

	for i := 0; i < len(tasks); i += workersCount {
		wg.Add(workersCount)

		for j := 0; j < workersCount; j++ {
			go func(taskIdx int) {
				defer wg.Done()
				errs <- tasks[taskIdx]()
			}(i + j)
		}

		wg.Wait()

		for j := 0; j < workersCount; j++ {
			if err := <-errs; err != nil {
				errsCount++
			}
		}

		if m > 0 && errsCount >= m {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
