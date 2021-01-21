package jobs

type JobManager interface {
	AddJobStateListener(func(JobState))
	Submit(func() ([]KV, error)) (JobId, error)
}

type simpleJobManager struct {
	runner            JobRunner
	jobStateListeners []func(JobState)
}

func NewSimpleJobManager(jobRunner JobRunner) JobManager {
	simple := &simpleJobManager{
		runner: jobRunner,
	}

	return simple
}

func (sqjm *simpleJobManager) AddJobStateListener(f func(JobState)) {
	sqjm.jobStateListeners = append(sqjm.jobStateListeners, f)
}

func (sqjm *simpleJobManager) Submit(jobFunc func() ([]KV, error)) (JobId, error) {
	job := NewJob(jobFunc, WithStateHandler(sqjm.handleState))
	job.stateFunc(submitted(job.id, "ok", job.initialMetadata))
	sqjm.runner.Run(job)
	return job.id, nil
}

func (sqjm *simpleJobManager) handleState(jobState JobState) {
	if sqjm.jobStateListeners != nil {
		for _, f := range sqjm.jobStateListeners {
			f(jobState)
		}
	}
}
