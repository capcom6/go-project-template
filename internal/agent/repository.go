package agent

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Repository struct {
	mu     sync.RWMutex
	items  []taskModel
	nextID int
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Add(ctx context.Context, prompt string) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	now := time.Now()
	model := taskModel{
		ID:        r.nextID,
		Status:    TaskStatusPending,
		Prompt:    prompt,
		CreatedAt: now,
	}
	r.items = append(r.items, model)

	return Task{
		ID:        model.ID,
		Status:    model.Status,
		Prompt:    model.Prompt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *Repository) List(ctx context.Context) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]Task, len(r.items))
	for i, item := range r.items {
		tasks[i] = r.toTask(item)
	}

	return tasks, nil
}

func (r *Repository) Get(ctx context.Context, id int) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.ID == id {
			return r.toTask(item), nil
		}
	}

	return Task{}, fmt.Errorf("get task %d: %w", id, ErrTaskNotFound)
}

func (r *Repository) NextPending(ctx context.Context) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, item := range r.items {
		if item.Status == TaskStatusPending {
			r.items[i].Status = TaskStatusProcessing
			return r.toTask(r.items[i]), nil
		}
	}

	return Task{}, ErrTaskNotFound
}

func (r *Repository) Update(ctx context.Context, task Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, item := range r.items {
		if item.ID == task.ID {
			r.items[i].Status = task.Status
			r.items[i].Result = task.Result
			r.items[i].Error = task.Error
			r.items[i].ProcessedAt = task.ProcessedAt
			return nil
		}
	}

	return fmt.Errorf("update task %d: %w", task.ID, ErrTaskNotFound)
}

func (r *Repository) PendingCount(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, item := range r.items {
		if item.Status == TaskStatusPending {
			count++
		}
	}

	return count
}

func (r *Repository) toTask(m taskModel) Task {
	return Task{
		ID:          m.ID,
		Status:      m.Status,
		Prompt:      m.Prompt,
		Result:      m.Result,
		Error:       m.Error,
		CreatedAt:   m.CreatedAt,
		ProcessedAt: m.ProcessedAt,
	}
}
