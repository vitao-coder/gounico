package workerPool

import (
	"context"
	"fmt"
	"gounico/pkg/worker"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan worker.WorkerJob
	results      chan worker.WorkerJobResult
	Done         chan struct{}
}

func NewWorkerPool(workersCount int) *WorkerPool {
	return &WorkerPool{
		workersCount: workersCount,
		jobs:         make(chan worker.WorkerJob, workersCount),
		results:      make(chan worker.WorkerJobResult, workersCount),
		Done:         make(chan struct{}),
	}
}

func start(ctx context.Context, wg *sync.WaitGroup, jobs <-chan worker.WorkerJob, results chan<- worker.WorkerJobResult) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {

				return
			}
			results <- job.ExecuteJob(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- worker.WorkerJobResult{
				Error: ctx.Err(),
			}
			return
		}
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go start(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) AddJobs(workJobs ...worker.WorkerJob) {
	for i := range workJobs {
		wp.jobs <- workJobs[i]
	}
}

func (wp *WorkerPool) Results() <-chan worker.WorkerJobResult {
	return wp.results
}

func (wp *WorkerPool) Finished() <-chan struct{} {
	return wp.Done
}
