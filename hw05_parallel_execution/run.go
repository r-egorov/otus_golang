package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func RunSeq(tasks []Task, n, m int) error {
	errsCount := 0
	for i := 0; i < len(tasks); i += n {
		if errsCount == m {
			return ErrErrorsLimitExceeded
		}
		err := tasks[i]()
		if err != nil {
			errsCount++
		}
	}
	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errsCount := 0

	routinesCount := len(tasks)
	if n < routinesCount {
		routinesCount = n
	}

	for i := 0; i < len(tasks); i += routinesCount {
		errs := make(chan error, routinesCount)
		wg.Add(routinesCount)

		for j := 0; j < routinesCount; j++ {
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
