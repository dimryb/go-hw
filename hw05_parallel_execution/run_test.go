package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
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

	t.Run("m equals 0, any error should return ErrErrorsLimitExceeded", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(10))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := 0

		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "expected ErrErrorsLimitExceeded, got %v", err)
		require.LessOrEqual(t, runTasksCount, int32(0), "too many tasks were started")
	})

	t.Run("m negative, ignore all errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(10))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 5
		maxErrorsCount := -1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err, "expected no error when m < 0")
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}

func TestTasksWithoutSleep(t *testing.T) {
	t.Run("tasks without errors using concurrency check", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		var wg sync.WaitGroup

		for i := 0; i < tasksCount; i++ {
			wg.Add(1)
			tasks = append(tasks, func() error {
				defer wg.Done()
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		errChan := make(chan error, 1)
		go func() {
			errChan <- Run(tasks, workersCount, maxErrorsCount)
		}()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		duration := 1 * time.Second
		require.Eventually(t, func() bool {
			currentCount := atomic.LoadInt32(&runTasksCount)
			if currentCount != int32(tasksCount) {
				t.Logf("Current runTasksCount: %d, expected: %d", currentCount, tasksCount)
			}
			return atomic.LoadInt32(&runTasksCount) == int32(tasksCount)
		}, duration, 10*time.Millisecond, "not all tasks were completed concurrently")

		select {
		case <-time.After(duration):
			t.Fatalf("Run did not complete within the timeout")
		case err := <-errChan:
			require.NoError(t, err, "expected no error during task execution")
		}

		select {
		case <-done:
			require.Equal(t, int32(tasksCount), atomic.LoadInt32(&runTasksCount), "not all tasks were completed")
		case <-time.After(duration):
			t.Fatalf("Some tasks did not complete within the timeout")
		}
	})
}
