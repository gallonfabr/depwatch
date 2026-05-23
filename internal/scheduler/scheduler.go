// Package scheduler provides a simple ticker-based job scheduler
// that triggers dependency polling at a configured interval.
package scheduler

import (
	"context"
	"log"
	"time"
)

// Job is a function that will be called on each tick.
type Job func(ctx context.Context) error

// Scheduler runs a Job at a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *log.Logger
}

// New creates a new Scheduler with the given interval and job.
func New(interval time.Duration, job Job, logger *log.Logger) (*Scheduler, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("scheduler: interval must be positive, got %s", interval)
	}
	if job == nil {
		return nil, fmt.Errorf("scheduler: job must not be nil")
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}, nil
}

// Run starts the scheduler loop, executing the job immediately and then
// at every interval. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)

	// Run once immediately before waiting for the first tick.
	s.runJob(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runJob(ctx)
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) {
	if err := s.job(ctx); err != nil {
		s.logger.Printf("scheduler: job error: %v", err)
	}
}
