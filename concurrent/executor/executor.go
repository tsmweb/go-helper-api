/*
Package executor provides an implementation of the Executor to perform background
processing, limiting resource consumption when executing a collection of jobs.

Executor example:

	exe := executor.New(10)
	defer exe.Shutdown()

	for i := 0; i < 100; i++ {
		idx := i
		err := exe.Schedule(func(ctx context.Context) {
			select {
			case <-ctx.Done():
				log.Printf("[JOB] ID %d - stop \n", idx)
				return
			default:
				log.Printf("[JOB] ID %d \n", idx)
				time.Sleep(time.Millisecond * 100)
		})
		if err != nil {
			t.Log(err)
			break
		}
	}

*/
package executor

import (
	"context"
	"errors"
	"sync"
)

// Executor contains logic for goroutine execution.
type Executor struct {
	sema chan struct{}
	wg sync.WaitGroup

	shutdown chan struct{}
	mu sync.Mutex // guard terminated
	terminated bool

	ctx context.Context
	cancelFunc context.CancelFunc
}

// New returns an initialized goroutine executor.
//
// The size parameter determines the limit of the executor.
func New(size int) *Executor {
	if size <= 0 {
		size = 1
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Executor{
		sema: make(chan struct{}, size),
		shutdown: make(chan struct{}),
		terminated: false,
		ctx: ctx,
		cancelFunc: cancelFunc,
	}
}

// Schedule schedules the work to be performed on the executors.
func (e *Executor) Schedule(job func(ctx context.Context)) error {
	if e.IsTerminated() {
		return errors.New("executor terminated")
	}
	e.sema <- struct{}{}
	e.worker(job)
	return nil
}

// Shutdown closes the executor.
func (e *Executor) Shutdown() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.terminated {
		close(e.shutdown) // signal to shutdown workers and close the executor.
		e.cancelFunc() // tells an operation to abandon its work.
		e.terminated = true
		e.wg.Wait()
	}
}

// IsTerminated return if the executor is terminated.
func (e *Executor) IsTerminated() bool {
	select {
	case <-e.shutdown:
		return true
	default:
		return false
	}
}

func (e *Executor) worker(job func(ctx context.Context)) {
	done := make(chan struct{})
	e.wg.Add(1)

	go func() {
		defer func() {
			<-e.sema
			e.wg.Done()
			//log.Println("[EXECUTOR] done")
		}()

		done <- struct{}{}
		job(e.ctx)
	}()

	<-done
}
