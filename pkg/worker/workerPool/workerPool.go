package workerPool

import (
	"context"
	"fmt"
	"gounico/pkg/worker/domain"
)

type WorkerPool struct {
	workersCount int
	jobs         chan domain.WorkerJob
	results      chan domain.WorkerJobResult
	Done         chan struct{}
}

func NewWorkerPool(workersCount int) *WorkerPool {
	return &WorkerPool{
		workersCount: workersCount,
		jobs:         make(chan domain.WorkerJob, workersCount),
		results:      make(chan domain.WorkerJobResult, workersCount),
		Done:         make(chan struct{}),
	}
}

func start(ctx context.Context, jobs <-chan domain.WorkerJob, results chan<- domain.WorkerJobResult) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {

				return
			}
			go func() {
				results <- job.ExecuteJob()
			}()
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- domain.WorkerJobResult{
				Error: ctx.Err(),
			}
			return
		}
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	for i := 0; i < wp.workersCount; i++ {
		go start(ctx, wp.jobs, wp.results)
	}
}

func (wp *WorkerPool) AddJobs(workJobs ...domain.WorkerJob) {
	for i := range workJobs {
		wp.jobs <- workJobs[i]
	}
}

func (wp *WorkerPool) Results() <-chan domain.WorkerJobResult {
	return wp.results
}

func (wp *WorkerPool) Finished() <-chan struct{} {
	return wp.Done
}
