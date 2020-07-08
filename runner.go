package jobs

type JobRunner interface {
	Run(Job)
}

func NewBoundlessJobRunner() JobRunner {
	return boundlessJobRunner{}
}

type boundlessJobRunner struct {
}

func (bjr boundlessJobRunner) Run(job Job) {
	id := job.id
	go func(JobId) {
		var err error
		var metadata map[string]interface{}
		metadata = job.initialMetadata
		defer func() {
			errMsg := "ok"
			if err != nil {
				errMsg = err.Error()
			}
			job.stateFunc(finished(id, errMsg, metadata))
		}()
		job.stateFunc(started(id, "ok", metadata))
		returnedMetadata, err := job.jobFunc()
		if returnedMetadata != nil {
			for k, v := range returnedMetadata {
				metadata[k] = v
			}
		}
	}(id)
}
