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
	go func() {
		var err error
		var metadata map[string]interface{}
		defer func() {
			errMsg := "ok"
			if err != nil {
				errMsg = err.Error()
			}
			job.stateFunc(finished(errMsg,metadata))
		}()
		metadata, err = job.jobFunc()
	}()
}