package processor

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/kachan1208/something/processor/dao/data"
	"github.com/kachan1208/something/processor/dao/storage"
	"github.com/kachan1208/something/processor/model"
	"github.com/kachan1208/something/processor/stream/buffer"
	"github.com/kachan1208/something/processor/stream/decoder"
)

type PortProcessor struct {
	decoder       decoder.StreamDecoder
	dataStor      data.DataStorage
	jobs          sync.Map
	jobQueue      chan struct{}
	buf           buffer.Stream
	storageClient storage.StorageClient
}

var (
	WorkersCount = 10
	MaxBatchTime = time.Second * 5
	MaxBatchSize = 1024
)

func NewPortProcessor(
	decoder decoder.StreamDecoder,
	dataStor data.DataStorage,
	buf buffer.Stream,
	storageClient storage.StorageClient,
) *PortProcessor {
	p := PortProcessor{
		dataStor:      dataStor,
		decoder:       decoder,
		jobs:          sync.Map{},
		jobQueue:      make(chan struct{}, 1),
		buf:           buf,
		storageClient: storageClient,
	}

	return &p
}

func (p *PortProcessor) init(j *job) {
	j.setStatus(StatusRunning)
	for i := 0; i < WorkersCount; i++ {
		go p.process(j, j.ctx)
	}
}

func (p *PortProcessor) Process(filename string) (string, error) {
	in, err := p.dataStor.Load(filename)
	if err != nil {
		return "", err
	}

	jobID, job := newJob(in)
	p.jobs.Store(jobID, job)

	go func(data io.ReadCloser) {
		var port model.Port
		defer p.dataStor.Close(data)

		//Simple lock mechanism to make queue of jobs
		//yes, new goroutine requires memory of service that can be overtaken, it's better to use db or queue
		//(I don't want to add a queue service, it's just a test task)
		p.jobQueue <- struct{}{}
		p.init(job)
		_, err := p.decoder.Decode(in, &port, p.buf)
		if err != nil {
			job.setError(err)
		} else {
			job.setStatus(StatusDone)
		}

		job.finish()
		<-p.jobQueue
	}(in)

	return jobID, nil
}

func (p *PortProcessor) GetStatus(jobID string) (int32, uint64, error) {
	j, ok := p.jobs.Load(jobID)
	if !ok {
		return 0, 0, errors.New("can't find job")
	}

	job := j.(*job)
	return job.getStatus(), job.getCounter(), nil
}

func (p *PortProcessor) process(j *job, ctx context.Context) {
	timer := time.NewTimer(MaxBatchTime)
	batch := make([]*model.Port, 0, MaxBatchSize)
	isExit := false

	for {
		select {
		case msg, _ := <-p.buf.Get():
			port, _ := msg.(*model.Port)
			batch = append(batch, port)

			if len(batch) < MaxBatchSize {
				continue
			}
		case <-timer.C:
		case <-ctx.Done():
			isExit = true
		}

		timer.Stop()
		timer.Reset(MaxBatchTime)

		if len(batch) > 0 {
			//circuit breaker should exist
			count, err := p.storageClient.StoreBatchOfPorts(batch)
			if err != nil {
				j.setError(err)
				return
			}

			j.updateCounter(count)
		}

		if isExit {
			return
		}
	}
}
