package executor

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestExecutor_Schedule(t *testing.T) {
	jobSize := 100

	exe := New(10)
	defer exe.Shutdown()

	for i := 0; i < jobSize; i++ {
		job := &PrintJob{
			Index:    i,
			Duration: time.Millisecond * 100,
		}

		err := exe.Schedule(job.Run())
		if err != nil {
			t.Log(err)
			break
		}
	}
}

func TestExecutor_Shutdown(t *testing.T) {
	jobSize := 20

	exe := New(2)
	defer exe.Shutdown()

	tm := time.NewTimer(time.Second * 1)
	go func() {
		<-tm.C
		exe.Shutdown()
	}()

	for i := 0; i < jobSize; i++ {
		job := &PrintJob{
			Index:    i,
			Duration: time.Millisecond * 500,
		}

		err := exe.Schedule(job.Run())
		if err != nil {
			t.Log(err)
			break
		}
	}
}

type PrintJob struct {
	Index    int
	Duration time.Duration
}

func (p *PrintJob) Run() func(ctx context.Context) {
	return func(ctx context.Context) {
		select {
		case <-ctx.Done():
			log.Printf("[JOB] ID %d - stop \n", p.Index)
			return
		default:
			log.Printf("[JOB] ID %d \n", p.Index)
			time.Sleep(p.Duration)
		}
	}
}
