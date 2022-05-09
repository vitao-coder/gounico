package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	ctx := context.Background()

	maxWorkers := 6

	workerTest := NewWorkerPool(maxWorkers)

	go workerTest.Run(ctx)

	job1 := Job1()
	job2 := Job2()
	job3 := Job3()
	job4 := Job3()
	job5 := Job2()
	job6 := Job1()

	go workerTest.AddJobs(job1, job2, job3)

	go workerTest.AddJobs(job4, job5, job6)

	for {
		select {
		case r, ok := <-workerTest.Results():
			if !ok {
				continue
			}
			t.Log(r.Result)
			t.Log(r.WorkerJobDescriptor)
			t.Log(r.Error)

		case <-workerTest.Done:
			return
		default:
		}
	}

}

func TestCancelWorkerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	maxWorkers := 5

	workerTest := NewWorkerPool(maxWorkers)

	job1 := Job1()
	job2 := Job2()
	job3 := Job3()
	job4 := Job3()
	job5 := Job2()
	job6 := Job1()

	go workerTest.AddJobs(job1, job2, job3, job4, job5, job6)
	go workerTest.Run(ctx)

	go func() {
		for {
			select {
			case r, ok := <-workerTest.Results():
				if !ok {
					continue
				}
				t.Log(r.Result)
				t.Log(r.WorkerJobDescriptor)
				t.Log(r.Error)

			case <-workerTest.Done:
				return
			default:
			}
		}
	}()
	time.Sleep(2 * time.Second)

	cancel()
}

func Job1() WorkerJob {
	var paramsTest []interface{}
	paramsTest = append(paramsTest, "paramA")
	paramsTest = append(paramsTest, "paramB")
	return NewWorkerJob("Job1", execTestFunction, paramsTest)
}

func Job2() WorkerJob {
	var paramsTest []interface{}
	paramsTest = append(paramsTest, "paramC")
	return NewWorkerJob("Job2", execTestFunction, paramsTest)
}

func Job3() WorkerJob {
	var paramsTest []interface{}
	paramsTest = append(paramsTest, "paramD")
	return NewWorkerJob("Job3", execTestFunctionError, paramsTest)
}

var (
	execTestFunction = func(ctx context.Context, params ...interface{}) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return params, nil
	}

	execTestFunctionError = func(ctx context.Context, params ...interface{}) (interface{}, error) {
		time.Sleep(2 * time.Second)
		return nil, errDefault
	}

	errDefault = errors.New("err")
)
