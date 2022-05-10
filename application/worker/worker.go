package worker

import (
	"context"
	"gounico/application/worker/domain"
)

type Worker interface {
	Run(ctx context.Context)
	AddJobs(workJobs ...domain.WorkerJob)
	Results() <-chan domain.WorkerJobResult
}
