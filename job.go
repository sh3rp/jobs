package jobs

type job struct {
	id              JobId
	jobFunc         func() (map[string]interface{}, error)
	stateFunc       func(JobState)
	initialMetadata map[string]interface{}
}

type JobId string

type JobState struct {
	id       JobId
	state    int
	message  string
	metadata map[string]interface{}
}

func (js JobState) Message() string {
	return js.message
}

func (js JobState) Metadata() map[string]interface{} {
	return js.metadata
}

func (js JobState) String() string {
	return [...]string{"submitted", "started", "finished"}[js.state]
}

const (
	SUBMITTED = iota
	STARTED
	FINISHED
)

func submitted(id JobId, msg string, metadata map[string]interface{}) JobState {
	return newJobState(id, SUBMITTED, msg, metadata)
}
func started(id JobId, msg string, metadata map[string]interface{}) JobState {
	return newJobState(id, STARTED, msg, metadata)
}

func finished(id JobId, msg string, metadata map[string]interface{}) JobState {
	return newJobState(id, FINISHED, msg, metadata)
}

func newJobState(id JobId, state int, msg string, metadata map[string]interface{}) JobState {
	return JobState{
		id:       id,
		state:    state,
		message:  msg,
		metadata: metadata,
	}
}

func WithStateHandler(h func(JobState)) func(*job) {
	return func(job *job) {
		job.stateFunc = h
	}
}

func WithKV(key string, value interface{}) func(*job) {
	return func(job *job) {
		if job.initialMetadata == nil {
			job.initialMetadata = make(map[string]interface{})
		}
		job.initialMetadata[key] = value
	}
}

func newJob(f func() (map[string]interface{}, error), options ...func(*job)) job {
	job := &job{
		id:      JobId(NewID()),
		jobFunc: f,
	}

	for _, opt := range options {
		opt(job)
	}

	return *job
}
