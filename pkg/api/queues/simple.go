// Handles job queues
package queues

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/rayaman/simply-jobber/pkg/models"
)

type Simple struct {
	queue      *queue.Queue
	processJob models.ProcessFunc
	mutex      sync.Mutex
	mutexq     sync.Mutex
	jid        int
	responses  []models.ResponseFunc
	handles    map[models.JobType]models.ProcessFunc
	workers    uint
	context    *context.Context
	cancel     context.CancelFunc
}

type pkg struct {
	models.Job
	jid int
}

func NewSimple(workers uint) *Simple {
	sq := &Simple{workers: workers, handles: map[models.JobType]models.ProcessFunc{}, queue: queue.New(), mutex: sync.Mutex{}, mutexq: sync.Mutex{}, processJob: func(j models.Job) (models.JobResponse, error) { return models.JobResponse{}, nil }}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered!")
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	sq.context = &ctx
	sq.cancel = cancel
	sq.genWorkers(workers)
	return sq
}

func (sq *Simple) Len() int {
	sq.mutexq.Lock()
	defer sq.mutexq.Unlock()
	return sq.queue.Len()
}

func (sq *Simple) genWorkers(workers uint) {
	for range workers {
		go func() {
			ctx := sq.context
			for {
				time.Sleep(time.Millisecond)
				select {
				case <-(*ctx).Done():
					return
				default:
					if p, ok := sq.dequeue(); ok {
						var err error
						var res models.JobResponse
						sq.mutex.Lock()
						if handle, ok := sq.handles[p.Type]; ok {
							sq.mutex.Unlock()
							res, err = handle(p.Job)
						} else {
							sq.mutex.Unlock()
							res, err = sq.processJob(p.Job)
						}
						if err == nil {
							for _, jr := range sq.responses {
								jr(res, p.jid)
							}
						} else {
							// Todo log
						}
					}
				}
			}
		}()
	}
}

func (sq *Simple) Scale(workers uint) {
	sq.mutex.Lock()
	defer sq.mutex.Unlock()
	sq.cancel()
	ctx, cancel := context.WithCancel(context.Background())
	sq.context = &ctx
	sq.cancel = cancel
	sq.genWorkers(workers)
}

func (sq *Simple) dequeue() (pkg, bool) {
	sq.mutexq.Lock()
	defer sq.mutexq.Unlock()
	if sq.queue.Len() > 0 {
		return sq.queue.Dequeue().(pkg), true
	} else {
		return pkg{}, false
	}
}

func (sq *Simple) SetProcessor(fn models.ProcessFunc, args ...string) error {
	sq.mutex.Lock()
	defer sq.mutex.Unlock()
	sq.processJob = fn
	return nil
}

func (sq *Simple) OnJobResponse(jr models.ResponseFunc) {
	sq.mutex.Lock()
	defer sq.mutex.Unlock()
	sq.responses = append(sq.responses, jr)
}

func (sq *Simple) AddHandle(t models.JobType, fn models.ProcessFunc) error {
	sq.mutex.Lock()
	defer sq.mutex.Unlock()
	if _, ok := sq.handles[t]; ok {
		return fmt.Errorf("job type [%v] already exists", t)
	}
	sq.handles[t] = fn
	return nil
}

func (sq *Simple) Send(j models.Job) (int, error) {
	sq.mutexq.Lock()
	defer sq.mutexq.Unlock()
	sq.queue.Enqueue(pkg{Job: j, jid: sq.jid})
	sq.jid++
	return sq.jid, nil
}
