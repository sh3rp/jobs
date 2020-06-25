package jobs

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestJobManagerOneJob(t *testing.T) {
	runner := NewBoundlessJobRunner()
	manager := NewSimpleJobManager(runner)
	wg := sync.WaitGroup{}
	wg.Add(1)
	manager.Submit(Job{
		id: JobId(NewID()),
		jobFunc: func() (map[string]interface{},error) {
			t.Logf("executing some bullshit")
			time.Sleep(1 * time.Second)
			t.Logf("doing some more bullshit")
			time.Sleep(1 * time.Second)
			return map[string]interface{}{
				"donkey":"fucker",
				"device":"fuck-ebay01",
			}, nil
		},
		stateFunc: func(state JobState) {
			t.Logf("[STATE] [%s] msg=%s, metadata=%+v",state.String(),state.Message(),state.Metadata())
			if state.state == FINISHED {
				wg.Done()
			}
		},
	})
	wg.Wait()
}

func TestJobManagerFiftyJobs(t *testing.T) {
	runner := NewBoundlessJobRunner()
	manager := NewSimpleJobManager(runner)
	wg := sync.WaitGroup{}
	for i := 0;i < 50;i++ {
		wg.Add(1)
		job := func() (map[string]interface{}, error) {
			t.Logf("executing some bullshit")
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			t.Logf("doing some more bullshit")
			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			return map[string]interface{}{
				"donkey": "fucker",
				"device": "fuck-ebay01",
			}, nil
		}
		manager.Submit(NewJob(job, WithStateHandler(func(state JobState) {
			t.Logf("[STATE] [%s] msg=%s, metadata=%+v", state.String(), state.Message(), state.Metadata())
			if state.state == FINISHED {
				wg.Done()
			}
		})))
	}
	wg.Wait()
}
