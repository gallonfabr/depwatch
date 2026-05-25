// Package scheduler provides a ticker-based job scheduler for depwatch.
//
// It accepts a Job function and a polling interval, running the job
// immediately on start and then repeatedly at the given interval until
// the provided context is cancelled.
//
// The scheduler guarantees that:
//   - The job is executed once immediately upon calling Run.
//   - Subsequent executions occur at the specified interval.
//   - Execution stops cleanly when the context is cancelled.
//   - A minimum interval of 1 second is enforced to prevent busy-looping.
//
// Typical usage:
//
//	s, err := scheduler.New(cfg.Interval, watcher.Poll, logger)
//	if err != nil {
//		log.Fatal(err)
//	}
//	s.Run(ctx)
package scheduler
