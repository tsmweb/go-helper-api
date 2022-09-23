package executor

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestExecutor_Schedule(t *testing.T) {
	taskSize := 100

	exe := New(10)
	defer exe.Shutdown()

	for i := 0; i < taskSize; i++ {
		task := &PrintTask{
			Index:    i,
			Duration: time.Millisecond * 100,
		}

		err := exe.Schedule(task.Run())
		if err != nil {
			t.Log(err)
			break
		}
	}
}

func TestExecutor_Shutdown(t *testing.T) {
	taskSize := 20

	exe := New(2)
	defer exe.Shutdown()

	tm := time.NewTimer(time.Second * 1)
	go func() {
		<-tm.C
		exe.Shutdown()
	}()

	for i := 0; i < taskSize; i++ {
		task := &PrintTask{
			Index:    i,
			Duration: time.Millisecond * 500,
		}

		err := exe.Schedule(task.Run())
		if err != nil {
			t.Log(err)
			break
		}
	}
}

type PrintTask struct {
	Index    int
	Duration time.Duration
}

func (p *PrintTask) Run() func(ctx context.Context) {
	return func(ctx context.Context) {
		select {
		case <-ctx.Done():
			log.Printf("[TASK] ID %d - stop \n", p.Index)
			return
		default:
			log.Printf("[TASK] ID %d \n", p.Index)
			time.Sleep(p.Duration)
		}
	}
}
