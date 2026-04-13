package sync

import (
	"context"
	"sync"
	"time"

	"devops1/internal/config"
	"devops1/internal/db"
	"devops1/internal/metrics"
	"devops1/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

type Engine struct {
	Cfg     *config.Config
	Queries *db.Queries
	Logger  zerolog.Logger
	Worker  *Worker
	LastRun time.Time
	Running bool
	mu      sync.RWMutex
}

func NewEngine(cfg *config.Config, q *db.Queries, logger zerolog.Logger) *Engine {
	return &Engine{Cfg: cfg, Queries: q, Logger: logger, Worker: &Worker{Cfg: cfg, Queries: q}}
}

func (e *Engine) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.Cfg.Sync.IntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = e.RunOnce(ctx)
		}
	}
}

func (e *Engine) RunOnce(ctx context.Context) error {
	e.mu.Lock()
	if e.Running {
		e.mu.Unlock()
		return nil
	}
	e.Running = true
	e.mu.Unlock()
	defer func() { e.mu.Lock(); e.Running = false; e.LastRun = time.Now(); e.mu.Unlock() }()
	start := time.Now()
	mappings, err := e.Queries.GetActiveMappings(ctx)
	if err != nil {
		return err
	}
	sem := semaphore.NewWeighted(int64(e.Cfg.Sync.Workers))
	wg := sync.WaitGroup{}
	for _, m := range mappings {
		m := m
		if err := sem.Acquire(ctx, 1); err != nil {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer sem.Release(1)
			if err := e.Worker.Sync(ctx, m); err != nil {
				_ = e.Queries.WriteLog(ctx, models.Log{Username: m.ExternalUser, Filename: "", Direction: "sync", Status: models.LogStatusFailed, Error: err.Error(), MappingID: m.ID, Component: "sync"})
			}
		}()
	}
	wg.Wait()
	metrics.SyncDurationSeconds.Observe(time.Since(start).Seconds())
	return nil
}

func (e *Engine) Status() (time.Time, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.LastRun, e.Running
}
