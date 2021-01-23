package jobs

func NewPooledJobRunner(poolSize int) JobRunner {
	return pooledJobRunner{poolSize, make(chan job), make(chan result)}
}

type pooledJobRunner struct {
	poolSize      int
	jobsChannel   chan job
	resultChannel chan result
}

type result struct {
	job job
	kvs []KV
	err error
}

func (pjr pooledJobRunner) start() {
	for i := 1; i <= pjr.poolSize; i++ {
		go worker(i, pjr.jobsChannel, pjr.resultChannel)
	}
}

func worker(workerId int, jobs chan job, results chan result) {
	for j := range jobs {
		var err error
		var metadata []KV
		metadata = j.initialMetadata
		defer func() {
			errMsg := "ok"
			if err != nil {
				errMsg = err.Error()
			}
			j.stateFunc(finished(j.id, errMsg, metadata))
		}()
		j.stateFunc(started(j.id, "ok", metadata))
		returnedMetadata, err := j.jobFunc()
		if returnedMetadata != nil && len(returnedMetadata) > 0 {
			for _, row := range returnedMetadata {
				metadata = append(metadata, row)
			}
		}
		results <- result{j, metadata, err}
	}
}

func (pjr pooledJobRunner) Run(j job) {
	pjr.jobsChannel <- j
}

func (pjr pooledJobRunner) publisher() {
	for result := range pjr.resultChannel {
		result.job.stateFunc(finished(job.id, result, metadata))
	}
}
