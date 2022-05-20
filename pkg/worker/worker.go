package worker

import (
	"context"
	"gounico/pkg/worker/domain"
)

type Worker interface {
	Run(ctx context.Context)
	AddJobs(workJobs ...domain.WorkerJob)
	Results() <-chan domain.WorkerJobResult
	Finished() <-chan struct{}
}
