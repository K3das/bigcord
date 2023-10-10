package jobs

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

type JobState int

const (
	JobStateInvalid JobState = iota
	JobStateRunning
	JobStateFailed
	JobStateCanceled
)

type Job struct {
	cancel context.CancelFunc
	State  JobState
	Error  error
}
type JobStore struct {
	log *zap.SugaredLogger

	jobsMu sync.RWMutex
	jobs   map[string]Job

	wg sync.WaitGroup

	ctx context.Context
}

func NewJobStore(ctx context.Context, rawLog *zap.Logger) *JobStore {
	return &JobStore{
		jobsMu: sync.RWMutex{},
		jobs:   make(map[string]Job),
		ctx:    ctx,
		log:    rawLog.Sugar().With("source", "job_pool"),
		wg:     sync.WaitGroup{},
	}
}

func (j *JobStore) createJob(id string, cancel context.CancelFunc) {
	j.jobsMu.Lock()
	j.jobs[id] = Job{
		cancel: cancel,
		State:  JobStateRunning,
	}
	j.jobsMu.Unlock()
}
func (j *JobStore) deleteJob(id string) {
	j.jobsMu.Lock()
	delete(j.jobs, id)
	j.jobsMu.Unlock()
}

func (j *JobStore) Add(id string, jobFunc func(ctx context.Context) error) error {
	j.jobsMu.RLock()
	job, ok := j.jobs[id]
	j.jobsMu.RUnlock()
	if ok && job.State == JobStateRunning {
		return fmt.Errorf("job already exists and running")
	}

	ctx, cancel := context.WithCancel(j.ctx)

	j.createJob(id, cancel)
	go func() {
		j.wg.Add(1)

		defer func() {
			cancel()
			j.wg.Done()
		}()

		err := jobFunc(ctx)
		if err == nil {
			j.deleteJob(id)
			return
		}

		jobErr := fmt.Errorf("error during job execution: %w", err)

		jobState := JobStateFailed
		if errors.Is(err, context.Canceled) {
			jobState = JobStateCanceled
		}

		j.jobsMu.Lock()
		j.jobs[id] = Job{
			cancel: nil,
			State:  jobState,
			Error:  jobErr,
		}
		j.jobsMu.Unlock()

		j.log.With("job_id", id).Error(jobErr)
	}()

	return nil
}

func (j *JobStore) CancelJob(id string) error {
	j.jobsMu.RLock()
	job, ok := j.jobs[id]
	j.jobsMu.RUnlock()
	if !ok || job.State != JobStateRunning {
		return fmt.Errorf("cannot cancel job that isn't running")
	}
	job.cancel()
	return nil
}

func (j *JobStore) RemoveFloating(id string) error {
	j.jobsMu.RLock()
	job, ok := j.jobs[id]
	j.jobsMu.RUnlock()
	if !ok || job.State == JobStateRunning {
		return fmt.Errorf("cannot remove job that isn't floating")
	}

	j.deleteJob(id)
	return nil
}

func (j *JobStore) RemoveAllFloating() {
	j.jobsMu.RLock()
	jobs := j.jobs
	j.jobsMu.RUnlock()

	for id, job := range jobs {
		if job.State == JobStateRunning {
			continue
		}
		j.deleteJob(id)
	}
}

func (j *JobStore) GetJob(id string) (*Job, error) {
	j.jobsMu.RLock()
	job, ok := j.jobs[id]
	j.jobsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("job does not exist")
	}
	job.cancel = nil
	return &job, nil
}

func (j *JobStore) GetJobs() map[string]Job {
	jobs := make(map[string]Job)
	j.jobsMu.RLock()
	for k, v := range j.jobs {
		jobs[k] = v
	}
	j.jobsMu.RUnlock()
	return jobs
}

func (j *JobStore) Wait() {
	j.wg.Wait()
}
