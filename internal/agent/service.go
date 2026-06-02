package agent

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	openrouter "github.com/revrost/go-openrouter"
	"go.uber.org/zap"
)

type Service struct {
	config  Config
	llm     *openrouter.Client
	repo    *Repository
	metrics *Metrics
	logger  *zap.Logger

	startedAt      atomic.Value
	processedCount atomic.Int64
}

func New(config Config, llm *openrouter.Client, repo *Repository, metrics *Metrics, logger *zap.Logger) *Service {
	s := &Service{
		config:  config,
		llm:     llm,
		repo:    repo,
		metrics: metrics,
		logger:  logger,
	}
	s.startedAt.Store(time.Now())

	return s
}

func (s *Service) Run(ctx context.Context) error {
	s.logger.Info("ai agent loop started",
		zap.String("model", s.config.Model),
		zap.Duration("poll_interval", s.config.PollInterval),
		zap.Int("batch_size", s.config.BatchSize),
	)

	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("ai agent loop stopped")
			return nil
		case <-ticker.C:
			s.processBatch(ctx)
		}
	}
}

func (s *Service) processBatch(ctx context.Context) {
	s.metrics.IncLoopIterations()

	for i := 0; i < s.config.BatchSize; i++ {
		task, err := s.repo.NextPending(ctx)
		if err != nil {
			return
		}

		s.processTask(ctx, task)
	}
}

func (s *Service) processTask(ctx context.Context, task Task) {
	s.logger.Info("processing task", zap.Int("task_id", task.ID), zap.String("prompt", task.Prompt))

	result, err := s.callLLM(ctx, task.Prompt)
	if err != nil {
		now := time.Now()
		task.Status = TaskStatusFailed
		task.Error = err.Error()
		task.ProcessedAt = &now

		if updateErr := s.repo.Update(ctx, task); updateErr != nil {
			s.logger.Error("failed to update failed task", zap.Int("task_id", task.ID), zap.Error(updateErr))
		}
		s.metrics.IncTasksProcessed(string(TaskStatusFailed))
		return
	}

	now := time.Now()
	task.Status = TaskStatusDone
	task.Result = result
	task.ProcessedAt = &now

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("failed to update completed task", zap.Int("task_id", task.ID), zap.Error(err))
	}
	s.processedCount.Add(1)
	s.metrics.IncTasksProcessed(string(TaskStatusDone))
	s.metrics.SetQueueDepth(s.repo.PendingCount(ctx))

	s.logger.Info("task processed",
		zap.Int("task_id", task.ID),
		zap.Int("result_length", len(result)),
	)
}

func (s *Service) callLLM(ctx context.Context, prompt string) (string, error) {
	s.metrics.IncLLMCalls()

	start := time.Now()
	defer func() {
		s.metrics.ObserveLLMDuration(time.Since(start).Seconds())
	}()

	resp, err := s.llm.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: s.config.Model,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleSystem,
				Content: openrouter.Content{Text: s.config.SystemPrompt},
			},
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: prompt},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("llm call: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("llm returned no choices")
	}

	return resp.Choices[0].Message.Content.Text, nil
}

type Status struct {
	Uptime         time.Duration `json:"uptime"`
	ProcessedCount int64         `json:"processed_count"`
	QueueDepth     int           `json:"queue_depth"`
	Model          string        `json:"model"`
}

func (s *Service) Status() Status {
	return Status{
		Uptime:         time.Since(s.startedAt.Load().(time.Time)),
		ProcessedCount: s.processedCount.Load(),
		QueueDepth:     s.repo.PendingCount(context.Background()),
		Model:          s.config.Model,
	}
}

func (s *Service) Tasks() ([]Task, error) {
	return s.repo.List(context.Background())
}

func (s *Service) Enqueue(ctx context.Context, prompt string) (Task, error) {
	task, err := s.repo.Add(ctx, prompt)
	if err != nil {
		return Task{}, fmt.Errorf("enqueue task: %w", err)
	}

	s.metrics.SetQueueDepth(s.repo.PendingCount(ctx))
	return task, nil
}
