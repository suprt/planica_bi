package cron

// Scheduler handles scheduled tasks
type Scheduler struct {
	// TODO: add sync service
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{}
}

// StartDailySync starts daily synchronization task
// Runs every night to update current month (M) data
func (s *Scheduler) StartDailySync() {
	// TODO: implement
	// Schedule: daily at night
}

// StartMonthlyFinalization starts monthly finalization task
// Runs on 1st of each month at 07:00 MSK
func (s *Scheduler) StartMonthlyFinalization() {
	// TODO: implement
	// Schedule: 1st day of month at 07:00 MSK
	// Finalizes previous month, updates M/M-1/M-2 window
}

