package workerPool

import (
	"context"
	"errors"
	"fmt"
	"gounico/internal/worker/domain"
	"gounico/utils"
	"testing"
	"time"
)

func TestWorkerPoolInfinitely(t *testing.T) {
	defer utils.TimeTrack(time.Now(), "TestWorkerPoolInfinitely")
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)

	maxWorkers := 200
	jobVolume := 20000

	workerTest := NewWorkerPool(maxWorkers)

	go workerTest.Run(ctx)

	go addALotOfJobs(workerTest, jobVolume)

	contVolume := 0

	for {
		select {
		case r, ok := <-workerTest.Results():
			if !ok {
				continue
			}
			t.Log(r.Result)
			t.Log(r.WorkerJobDescriptor)
			if r.Error != nil {
				t.Log(r.Error.Error())
			}
			contVolume++
			if contVolume >= jobVolume {
				t.Logf("Worker pool executed successful...")
				close(workerTest.jobs)
				close(workerTest.results)
				close(workerTest.Done)
				ctx.Done()
			}
		case <-workerTest.Done:
			return
		default:
		}
	}
}

func TestCancelWorkerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	maxWorkers := 4

	workerTest := NewWorkerPool(maxWorkers)

	job1 := job("Job1")
	job2 := jobError("Job2")
	job3 := job("Job3")
	job4 := job("Job4")
	job5 := job("Job5")
	job6 := job("Job6")

	go workerTest.AddJobs(job1, job2, job3, job4, job5, job6)
	go workerTest.Run(ctx)

	for {
		select {
		case r, ok := <-workerTest.Results():
			if !ok {
				continue
			}
			t.Log(r.Result)
			t.Log(r.WorkerJobDescriptor)
			if r.Error != nil {
				cancel()
				t.Log(r.Error.Error())
				return
			}

		case <-workerTest.Done:
			return
		default:
		}
	}

}

func addALotOfJobs(wp *WorkerPool, jobVolume int) {
	for i := 0; i < jobVolume; i++ {
		jb := job(fmt.Sprintf("job %v", i))
		go wp.AddJobs(jb)
	}
}

func job(name string) domain.WorkerJob {
	var paramsTest []interface{}
	paramsTest = append(paramsTest, "paramA")
	paramsTest = append(paramsTest, "paramB")
	return domain.NewWorkerJob(name, execTestFunction, paramsTest)
}

func jobError(name string) domain.WorkerJob {
	var paramsTest []interface{}
	paramsTest = append(paramsTest, "paramD")
	return domain.NewWorkerJob(name, execTestFunctionError, paramsTest)
}

var (
	execTestFunction = func(ctx context.Context, params ...interface{}) (interface{}, error) {
		return params, nil
	}

	execTestFunctionError = func(ctx context.Context, params ...interface{}) (interface{}, error) {
		return nil, errDefault
	}

	errDefault = errors.New("err")
)
