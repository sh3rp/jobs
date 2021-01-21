package jobs

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobManagerOneJob(t *testing.T) {
	manager := NewSimpleJobManager(NewBoundlessJobRunner())
	wg := sync.WaitGroup{}
	jobFunc := func() ([]KV, error) {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		return []KV{
			KV{K: K("donkey"), V: V{Type: "string", Value: "shoes"}},
			KV{K: K("device"), V: V{Type: "string", Value: "sw1-iad01"}},
			KV{K: K("deviceVersion"), V: V{Type: "int", Value: 34}},
		}, nil
	}
	wg.Add(1)
	manager.AddJobStateListener(func(state JobState) {
		t.Logf("[STATE] (%s) [%s] msg=%s, metadata=%+v", state.id, state.String(), state.Message(), state.Metadata())
		if state.state == FINISHED {
			wg.Done()
		}
	})
	manager.Submit(jobFunc)
	wg.Wait()
}

func TestJobManagerFiftyJobs(t *testing.T) {
	manager := NewSimpleJobManager(NewBoundlessJobRunner())
	wg := sync.WaitGroup{}
	jobFunc := func() ([]KV, error) {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		return []KV{
			KV{K: K("donkey"), V: V{Type: "string", Value: "shoes"}},
			KV{K: K("device"), V: V{Type: "string", Value: "sw1-iad01"}},
		}, nil
	}
	jobsDone := make(map[JobId]bool)
	manager.AddJobStateListener(func(state JobState) {
		t.Logf("[STATE] (%s) [%s] msg=%s, metadata=%+v", state.id, state.String(), state.Message(), state.Metadata())
		if state.state == FINISHED {
			jobsDone[state.id] = true
			wg.Done()
		}
	})
	for i := 0; i < 50; i++ {
		wg.Add(1)
		jobId, err := manager.Submit(jobFunc)
		assert.NoError(t, err)
		assert.NotNil(t, jobId)
	}
	wg.Wait()

	for _, val := range jobsDone {
		assert.True(t, val)
	}
}
