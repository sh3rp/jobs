package jobs

type Job struct {
	id JobId
	jobFunc func() (map[string]interface{},error)
	stateFunc func(JobState)
}

type JobId string

type JobState struct {
	state int
	message string
	metadata map[string]interface{}
}

func (js JobState) Message() string {
	return js.message
}

func (js JobState) Metadata() map[string]interface{} {
	return js.metadata
}

func (js JobState) String() string {
	return [...]string{"submitted","started","finished"}[js.state]
}

const (
	SUBMITTED = iota
	STARTED
	FINISHED
)

func submitted(msg string, metadata map[string]interface{}) JobState {
	return newJobState(SUBMITTED,msg,metadata)
}
func started(msg string, metadata map[string]interface{}) JobState {
	return newJobState(STARTED,msg,metadata)
}

func finished(msg string, metadata map[string]interface{}) JobState {
	return newJobState(FINISHED,msg,metadata)
}

func newJobState(state int, msg string, metadata map[string]interface{}) JobState {
	return JobState{
		state: state,
		message: msg,
		metadata: metadata,
	}
}

func WithStateHandler(h func(JobState)) func(*Job) {
	return func(job *Job) {
		job.stateFunc = h
	}
}

func NewJob(f func() (map[string]interface{},error), options ...func(*Job)) Job {
	job := &Job{
		id: JobId(NewID()),
		jobFunc: f,
	}

	for _,opt := range options {
		opt(job)
	}

	return *job
}

