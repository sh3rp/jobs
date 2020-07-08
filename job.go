package jobs

import "reflect"

type job struct {
	id              JobId
	jobFunc         func() ([]KV, error)
	stateFunc       func(JobState)
	initialMetadata []KV
}

type JobId string

type JobState struct {
	id       JobId
	state    int
	message  string
	metadata []KV
}

func (js JobState) Message() string {
	return js.message
}

func (js JobState) Metadata() []KV {
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

func submitted(id JobId, msg string, metadata []KV) JobState {
	return newJobState(id, SUBMITTED, msg, metadata)
}
func started(id JobId, msg string, metadata []KV) JobState {
	return newJobState(id, STARTED, msg, metadata)
}

func finished(id JobId, msg string, metadata []KV) JobState {
	return newJobState(id, FINISHED, msg, metadata)
}

func newJobState(id JobId, state int, msg string, metadata []KV) JobState {
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
		job.initialMetadata = append(job.initialMetadata, KV{K: K(key), V: V{Type: reflect.TypeOf(value).Name(), Value: value}})
	}
}

func NewJob(f func() ([]KV, error), options ...func(*job)) job {
	job := &job{
		id:      JobId(NewID()),
		jobFunc: f,
	}

	for _, opt := range options {
		opt(job)
	}

	return *job
}
