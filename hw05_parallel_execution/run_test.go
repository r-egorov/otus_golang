package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

type testCase struct {
	tasksCount, workersCount, maxErrorsCount int
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		maxErrorsCount := 23

		// first m tasks return an error
		for i := 0; i < maxErrorsCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		// other tasks are cool
		for i := maxErrorsCount; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 10
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors using sleep", func(t *testing.T) {
		t.Skip()
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("Run with M <= 0", func(t *testing.T) {
		// In this case the Run function should ignore all the mistakes
		tasksCount := 50
		workersCount := 5
		maxErrorsCount := 0

		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		// Going to be 25 errors and 25 successes
		for i := 0; i < tasksCount/2; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("tasks without errors without sleep", func(t *testing.T) {
		cases := []testCase{
			{50, 5, 1},
			{5, 50, 1},
			{50, 50, 1},
		}
		for _, tc := range cases {
			testNoErrorsNoSleep(t, tc)
		}
	})
}

func testNoErrorsNoSleep(t *testing.T, tc testCase) {
	t.Helper()
	tasks := make([]Task, 0, tc.tasksCount)
	var runTasksCount int32

	block := make(chan bool)
	for i := 0; i < tc.tasksCount; i++ {
		tasks = append(tasks, func() error {
			<-block
			atomic.AddInt32(&runTasksCount, 1)
			return nil
		})
	}

	runDone := make(chan bool)
	go func() {
		err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
		require.NoError(t, err)
		runDone <- true
		close(runDone)
	}()

	close(block)

	isFinished := func() bool {
		select {
		case <-runDone:
			return true
		default:
			return false
		}
	}

	require.Eventually(t, isFinished, time.Second*1, time.Millisecond*1, "tasks were run sequentially?")
	require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
}
