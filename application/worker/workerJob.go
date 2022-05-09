package worker

import "context"

type JobFunction func(ctx context.Context, params ...interface{}) (interface{}, error)

type WorkerJobResult struct {
	Result              interface{}
	Error               error
	WorkerJobDescriptor string
}

type WorkerJob struct {
	Params     []interface{}
	Descriptor string
	Job        JobFunction
}

func NewWorkerJob(workerJobDescriptor string, jobFunc JobFunction, params ...interface{}) WorkerJob {
	return WorkerJob{
		Params:     params,
		Descriptor: workerJobDescriptor,
		Job:        jobFunc,
	}
}

func (wj WorkerJob) executeJob(ctx context.Context) WorkerJobResult {
	result, err := wj.Job(ctx, wj.Params)

	if err != nil {
		return WorkerJobResult{
			Error:               err,
			WorkerJobDescriptor: wj.Descriptor,
		}
	}

	return WorkerJobResult{
		Result:              result,
		WorkerJobDescriptor: wj.Descriptor,
	}
}
