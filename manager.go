package jobs

type JobManager interface {
	Submit(Job) error
}

type simpleQueueJobManager struct {
	jobQueue chan Job
	runner JobRunner
}

func NewSimpleQueueJobManager(queueSize int, jobRunner JobRunner) JobManager {
	simple := simpleQueueJobManager {
		jobQueue: make(chan Job,queueSize),
		runner: jobRunner,
	}

	go simple.dequeingLoop()
	return simple
}

func (sqjm simpleQueueJobManager) Submit(job Job) error {
	if job.stateFunc != nil {
		job.stateFunc(submitted("ok",nil))
	}
	sqjm.jobQueue <- job
	return nil
}

func (sqjm simpleQueueJobManager) dequeingLoop() {
	select {
	case job :=<-sqjm.jobQueue:
		if job.stateFunc != nil {
			job.stateFunc(started("ok",nil))
		}
		sqjm.runner.Run(job)
	}
}

