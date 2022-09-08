/*
Package metric implements routines to collect metrics from localhost and send to a topic
in Apache Kafka. The metrics collected are: "uptime", "os", "total memory", "memory used",
"cpu count", "cpu user", "cpu system", "cpu idle" and "num goroutines".

Starts collecting and sending metrics:

	producer := kafka.New([]string{"localhost:9094"}, "CLIENT_ID").NewProducer("TOPIC_NAME")
	err := metric.Start("localhost", 10, producer)
	if err != nil {
		// ...
	defer metric.Stop()
*/

package metric

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/tsmweb/go-helper-api/kafka"
)

// metric structure represents a metric log.
type metric struct {
	Host        string  `json:"host"`
	Uptime      string  `json:"uptime"`
	OS          string  `json:"os"`
	MemoryTotal uint64  `json:"memory_total"`
	MemoryUsed  uint64  `json:"memory_used"`
	CPUCount    int     `json:"cpu_count"`
	CPUUser     float64 `json:"cpu_user"`
	CPUSystem   float64 `json:"cpu_system"`
	CPUIdle     float64 `json:"cpu_idle"`
	Goroutines  int     `json:"goroutines"`
	Timestamp   string  `json:"timestamp"`
}

// newMetric creates a metric instance.
func newMetric(host string) (*metric, error) {
	_uptime, err := uptime.Get()
	if err != nil {
		return nil, err
	}

	mem, err := memory.Get()
	if err != nil {
		return nil, err
	}

	cpuBefore, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	cpuAfter, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	memTotal := uint64((mem.Total / 1024) / 1024) // MB
	memUsed := uint64((mem.Used / 1024) / 1024)   // MB

	cpuTotal := float64(cpuAfter.Total - cpuBefore.Total)
	cpuUser := float64(cpuAfter.User-cpuBefore.User) / cpuTotal * 100
	cpuSystem := float64(cpuAfter.System-cpuBefore.System) / cpuTotal * 100
	cpuIdle := float64(cpuAfter.Idle-cpuBefore.Idle) / cpuTotal * 100

	m := &metric{
		Host:        host,
		Uptime:      _uptime.String(),
		OS:          fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH),
		MemoryTotal: memTotal,
		MemoryUsed:  memUsed,
		CPUCount:    runtime.NumCPU(),
		CPUUser:     cpuUser,
		CPUSystem:   cpuSystem,
		CPUIdle:     cpuIdle,
		Goroutines:  runtime.NumGoroutine(),
		Timestamp:   time.Now().Format("2006-01-02T15:04:05-0700"), // yyyy-MM-dd'T'HH:mm:ssZ
	}

	return m, nil
}

func (m metric) toJSON() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

var (
	ticker  *time.Ticker
	done    chan bool
	wg      sync.WaitGroup
	running bool
	mu      sync.RWMutex // guard running

	ErrRunning = errors.New("is already running")
)

// Start initializes collecting and sending metric data from localhost to the Apache Kafka topic.
func Start(host string, seconds int, producer kafka.Producer) error {
	mu.RLock()
	if running {
		mu.RUnlock()
		return ErrRunning
	}
	mu.RUnlock()

	ctx := context.Background()
	ticker = time.NewTicker(time.Duration(seconds) * time.Second)
	done = make(chan bool)
	running = true
	wg.Add(1)

	go func() {
		defer wg.Done()

	loop:
		for {
			select {
			case <-done:
				break loop

			case <-ticker.C:
				m, err := newMetric(host)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					break
				}

				if err = producer.Publish(ctx, []byte(host), m.toJSON()); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
				}
			}
		}

		producer.Close()
	}()

	return nil
}

// Stop ends sending metrics to the Apache Kafka producer.
func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if running {
		running = false
		ticker.Stop()
		done <- true
		wg.Wait()
	}
}
