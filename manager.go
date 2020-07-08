package jobs

import "fmt"

type JobManager interface {
	Submit(Job) error
}

type simpleJobManager struct {
	runner JobRunner
}

func NewSimpleJobManager(jobRunner JobRunner) JobManager {
	simple := simpleJobManager{
		runner: jobRunner,
	}

	return simple
}

func (sqjm simpleJobManager) Submit(job Job) error {
	if job.id == "" {
		job.id = JobId(NewID())
	}
	if job.jobFunc == nil {
		return fmt.Errorf("Job function must be defined")
	}
	if job.stateFunc == nil {
		return fmt.Errorf("State function must be defined")
	}
	job.stateFunc(submitted(job.id, "ok", job.initialMetadata))
	sqjm.runner.Run(job)
	return nil
}
