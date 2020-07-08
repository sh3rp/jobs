package jobs

type JobRunner interface {
	Run(job)
}

func NewBoundlessJobRunner() JobRunner {
	return boundlessJobRunner{}
}

type boundlessJobRunner struct {
}

func (bjr boundlessJobRunner) Run(job job) {
	id := job.id
	go func(JobId) {
		var err error
		var metadata []KV
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
		if returnedMetadata != nil && len(returnedMetadata) > 0 {
			for _, row := range returnedMetadata {
				metadata = append(metadata, row)
			}
		}
	}(id)
}
