package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsCounter struct {
	errCount int32
	limit    int
	cancel   context.CancelFunc
}

func (e *ErrorsCounter) Increment() {
	newCount := atomic.AddInt32(&e.errCount, 1)
	if newCount >= int32(e.limit) {
		e.cancel()
	}
}

func (e *ErrorsCounter) isLimitExceeded() bool {
	return atomic.LoadInt32(&e.errCount) >= int32(e.limit)
}

func Worker(ctx context.Context, tasks <-chan Task, counter *ErrorsCounter, workerNum int) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-tasks:
			if !ok {
				return
			}
			if !counter.isLimitExceeded() {
				err := t()
				if err != nil {
					counter.Increment()
				}
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksChan := make(chan Task)
	var wgWorkers sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	errCounter := &ErrorsCounter{limit: m, cancel: cancel}

	for i := 0; i < n; i++ {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			Worker(ctx, tasksChan, errCounter, i)
		}()
	}

	for _, task := range tasks {
		if errCounter.isLimitExceeded() {
			break
		}
		tasksChan <- task
	}
	close(tasksChan)

	wgWorkers.Wait()

	if errCounter.isLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
