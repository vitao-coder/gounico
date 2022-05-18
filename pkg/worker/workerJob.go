package worker

import (
	"context"
	"fmt"
	"gounico/pkg/telemetry/openTelemetry"
)

type JobFunction func(ctx context.Context, params []interface{}) (interface{}, error)

type WorkerJobResult struct {
	Result              interface{}
	Error               error
	WorkerJobDescriptor string
}

type WorkerJob struct {
	Params     []interface{}
	Descriptor string
	Job        JobFunction
	ctx        context.Context
}

func NewWorkerJob(workerJobDescriptor string, jobFunc JobFunction, ctx context.Context, params ...interface{}) WorkerJob {
	return WorkerJob{
		Params:     params,
		Descriptor: workerJobDescriptor,
		Job:        jobFunc,
		ctx:        ctx,
	}
}

func (wj WorkerJob) ExecuteJob() WorkerJobResult {
	ctx, traceSpan := openTelemetry.NewSpan(wj.ctx, fmt.Sprintf("WorkerJob.ExecuteJob - %s", wj.Descriptor))
	defer traceSpan.End()

	result, err := wj.Job(ctx, wj.Params)
	if err != nil {
		openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
		openTelemetry.AddSpanError(traceSpan, err)
		return WorkerJobResult{
			Error:               err,
			WorkerJobDescriptor: wj.Descriptor,
		}
	}
	openTelemetry.SuccessSpan(traceSpan, "Success")
	return WorkerJobResult{
		Result:              result,
		WorkerJobDescriptor: wj.Descriptor,
	}
}
