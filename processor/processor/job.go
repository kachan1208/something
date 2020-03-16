package processor

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	uuid "github.com/satori/go.uuid"
)

const (
	StatusQueued int32 = iota
	StatusRunning
	StatusDone
	StatusError
)

type job struct {
	file    io.ReadCloser
	counter uint64
	status  int32
	ctx     context.Context
	cancel  context.CancelFunc
	err     error
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func newJob(file io.ReadCloser) (string, *job) {
	jobID := uuid.NewV4().String()
	ctx, cancel := context.WithCancel(context.Background())
	j := job{
		file:    file,
		counter: uint64(0),
		status:  StatusQueued,
		ctx:     ctx,
		cancel:  cancel,
	}

	return jobID, &j
}

func (j *job) setStatus(val int32) {
	atomic.StoreInt32(&j.status, val)
}

func (j *job) getStatus() int32 {
	return atomic.LoadInt32(&j.status)
}

func (j *job) updateCounter(val uint64) {
	atomic.AddUint64(&j.counter, val)
}

func (j *job) getCounter() uint64 {
	return atomic.LoadUint64(&j.counter)
}

func (j *job) setError(err error) {
	j.mutex.Lock()
	//save only first
	if j.err != nil {
		j.err = err
		j.setStatus(StatusError)
	}
	j.mutex.Unlock()
}

func (j *job) finish() {
	j.cancel()
}
