package worker

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan WorkerJob
	results      chan WorkerJobResult
	Done         chan struct{}
}

func NewWorkerPool(workersCount int) WorkerPool {
	return WorkerPool{
		workersCount: workersCount,
		jobs:         make(chan WorkerJob, workersCount),
		results:      make(chan WorkerJobResult, workersCount),
		Done:         make(chan struct{}),
	}
}

func start(ctx context.Context, wg *sync.WaitGroup, jobs <-chan WorkerJob, results chan<- WorkerJobResult) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {

				return
			}
			results <- job.executeJob(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- WorkerJobResult{
				Error: ctx.Err(),
			}
			return
		}
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go start(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp WorkerPool) AddJobs(workJobs ...WorkerJob) {
	for i := range workJobs {
		wp.jobs <- workJobs[i]
	}
}

func (wp WorkerPool) Results() <-chan WorkerJobResult {
	return wp.results
}
