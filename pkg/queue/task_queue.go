package queue

import (
	"context"
	"sync"
)

type Task interface {
	Execute(ctx context.Context) error
}

type TaskQueue struct {
	tasks    chan Task
	workers  int
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewTaskQueue(bufferSize int, workers int) *TaskQueue {
	return &TaskQueue{
		tasks:    make(chan Task, bufferSize),
		workers:  workers,
		stopChan: make(chan struct{}),
	}
}

func (q *TaskQueue) Start(ctx context.Context) {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(ctx)
	}
}

func (q *TaskQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
}

func (q *TaskQueue) AddTask(task Task) {
	q.tasks <- task
}

func (q *TaskQueue) worker(ctx context.Context) {
	defer q.wg.Done()

	for {
		select {
		case task := <-q.tasks:
			_ = task.Execute(ctx)
		case <-q.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}
