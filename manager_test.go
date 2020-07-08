package jobs

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestJobManagerOneJob(t *testing.T) {
	manager := NewSimpleJobManager(NewBoundlessJobRunner())
	wg := sync.WaitGroup{}
	job := func() ([]KV, error) {
		t.Logf("executing some bullshit")
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		t.Logf("doing some more bullshit")
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		return []KV{
			KV{K: K("donkey"), V: V{Type: "string", Value: "shoes"}},
			KV{K: K("device"), V: V{Type: "string", Value: "sw1-iad01"}},
			KV{K: K("deviceVersion"), V: V{Type: "int", Value: 34}},
		}, nil
	}
	wg.Add(1)
	manager.Submit(NewJob(job,
		WithStateHandler(func(state JobState) {
			t.Logf("[STATE] (%s) [%s] msg=%s, metadata=%+v", state.id, state.String(), state.Message(), state.Metadata())
			if state.state == FINISHED {
				wg.Done()
			}
		}),
		WithKV("test", "iteration")))
	wg.Wait()
}

func TestJobManagerFiftyJobs(t *testing.T) {
	manager := NewSimpleJobManager(NewBoundlessJobRunner())
	wg := sync.WaitGroup{}
	job := func() ([]KV, error) {
		t.Logf("executing some bullshit")
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		t.Logf("doing some more bullshit")
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		return []KV{
			KV{K: K("donkey"), V: V{Type: "string", Value: "shoes"}},
			KV{K: K("device"), V: V{Type: "string", Value: "sw1-iad01"}},
		}, nil
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		manager.Submit(
			NewJob(job,
				WithStateHandler(func(state JobState) {
					t.Logf("[STATE] (%s) [%s] msg=%s, metadata=%+v", state.id, state.String(), state.Message(), state.Metadata())
					if state.state == FINISHED {
						wg.Done()
					}
				}),
				WithKV(fmt.Sprintf("test-%d", i), fmt.Sprintf("iteration-%d", i))))
	}
	wg.Wait()
}
