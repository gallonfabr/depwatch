package scheduler_test

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/depwatch/internal/scheduler"
)

var silentLogger = log.New(io.Discard, "", 0)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := scheduler.New(0, func(ctx context.Context) error { return nil }, silentLogger)
	if err == nil {
		t.Fatal("expected error for zero interval, got nil")
	}
}

func TestNew_NilJob(t *testing.T) {
	_, err := scheduler.New(time.Second, nil, silentLogger)
	if err == nil {
		t.Fatal("expected error for nil job, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := scheduler.New(time.Second, func(ctx context.Context) error { return nil }, silentLogger)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRun_ExecutesJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	s, _ := scheduler.New(10*time.Second, job, silentLogger)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to run at least once immediately")
	}
}

func TestRun_JobErrorDoesNotStop(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("boom")
	}

	s, _ := scheduler.New(20*time.Millisecond, job, silentLogger)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(80 * time.Millisecond)
	cancel()
	<-done

	if atomic.LoadInt32(&count) < 2 {
		t.Errorf("expected multiple runs despite errors, got %d", count)
	}
}
