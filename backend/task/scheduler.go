package task

import (
	"sync"
	"time"

	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/config"
	"github.com/PingWangWang/DeepSeek-Usage-WinGUI/backend/util"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	jobs   map[string]*ScheduledJob
	ticker *time.Ticker
	stopCh chan struct{}
	mu     sync.RWMutex
}

// ScheduledJob 定时任务
type ScheduledJob struct {
	ID       string
	Interval time.Duration
	LastRun  time.Time
	Callback func()
	Enabled  bool
}

// NewScheduler 创建新的调度器
func NewScheduler() *Scheduler {
	scheduler := &Scheduler{
		jobs:   make(map[string]*ScheduledJob),
		ticker: time.NewTicker(1 * time.Minute),
		stopCh: make(chan struct{}),
	}

	// 启动后台 ticker
	go scheduler.run()

	return scheduler
}

// RegisterAutoRefresh 注册自动刷新任务
func (s *Scheduler) RegisterAutoRefresh(interval time.Duration, callback func()) error {
	job := &ScheduledJob{
		ID:       "auto-refresh",
		Interval: interval,
		Callback: callback,
		Enabled:  true,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs["auto-refresh"] = job
	util.Logger.Info("Auto refresh job registered", "interval", interval)
	return nil
}

// RegisterSubscriptionCheck 注册订阅检查任务
func (s *Scheduler) RegisterSubscriptionCheck(sub *config.Subscription, callback func()) error {
	job := &ScheduledJob{
		ID:       "subscription-" + sub.ID,
		Interval: 5 * time.Minute, // 每 5 分钟检查一次
		Callback: callback,
		Enabled:  sub.Enabled,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs["subscription-"+sub.ID] = job
	util.Logger.Info("Subscription check job registered", "id", sub.ID)
	return nil
}

// run 后台运行循环
func (s *Scheduler) run() {
	for {
		select {
		case <-s.ticker.C:
			s.checkAndRun()
		case <-s.stopCh:
			return
		}
	}
}

// checkAndRun 检查并执行到期的任务
func (s *Scheduler) checkAndRun() {
	s.mu.RLock()
	jobs := make([]*ScheduledJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	s.mu.RUnlock()

	now := time.Now()
	for _, job := range jobs {
		if !job.Enabled {
			continue
		}

		// 检查是否应该运行
		if job.LastRun.IsZero() || now.Sub(job.LastRun) >= job.Interval {
			util.Logger.Debug("Executing scheduled job", "id", job.ID)
			job.Callback()
			job.LastRun = now
		}
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopCh)
	util.Logger.Info("Scheduler stopped")
}

// RemoveJob 移除任务
func (s *Scheduler) RemoveJob(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.jobs, id)
	util.Logger.Info("Job removed", "id", id)
}

// EnableJob 启用任务
func (s *Scheduler) EnableJob(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, ok := s.jobs[id]; ok {
		job.Enabled = true
		util.Logger.Info("Job enabled", "id", id)
	}
}

// DisableJob 禁用任务
func (s *Scheduler) DisableJob(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, ok := s.jobs[id]; ok {
		job.Enabled = false
		util.Logger.Info("Job disabled", "id", id)
	}
}
