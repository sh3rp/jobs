package jobs

type JobManager interface {
	Submit(Job) error
}

type simpleJobManager struct {
	runner   JobRunner
}

func NewSimpleJobManager(jobRunner JobRunner) JobManager {
	simple := simpleJobManager {
		runner: jobRunner,
	}

	return simple
}

func (sqjm simpleJobManager) Submit(job Job) error {
	if job.stateFunc != nil {
		job.stateFunc(submitted("ok",nil))
	}
	sqjm.runner.Run(job)
	return nil
}

