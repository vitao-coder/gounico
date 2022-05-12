package worker

import (
	"context"
)

type Worker interface {
	Run(ctx context.Context)
	AddJobs(workJobs ...WorkerJob)
	Results() <-chan WorkerJobResult
	Finished() <-chan struct{}
}
