package domain

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
	var resultWork WorkerJobResult
	go func() {
		var result interface{}
		var err error
		traceSpan := openTelemetry.SpanFromContext(wj.ctx)
		defer traceSpan.End()
		result, err = wj.Job(wj.ctx, wj.Params)
		if err != nil {
			openTelemetry.FailSpan(traceSpan, fmt.Sprintf("Error: %s", err.Error()))
			openTelemetry.AddSpanError(traceSpan, err)
			resultWork = WorkerJobResult{
				Error:               err,
				WorkerJobDescriptor: wj.Descriptor,
			}
			return
		}
		openTelemetry.SuccessSpan(traceSpan, fmt.Sprintf("Success"))
		resultWork = WorkerJobResult{
			Result:              result,
			WorkerJobDescriptor: wj.Descriptor,
		}
		return
	}()
	return resultWork
}
