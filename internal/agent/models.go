package agent

import "time"

type taskModel struct {
	ID          int
	Status      TaskStatus
	Prompt      string
	Result      string
	Error       string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
