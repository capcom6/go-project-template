package agent

import "time"

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID          int
	Status      TaskStatus
	Prompt      string
	Result      string
	Error       string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
