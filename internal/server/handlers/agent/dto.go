package agent

import "github.com/capcom6/go-project-template/internal/agent"

type EnqueueRequest struct {
	Prompt string `json:"prompt" validate:"required"`
}

type TaskResponse struct {
	ID          int     `json:"id"`
	Status      string  `json:"status"`
	Prompt      string  `json:"prompt"`
	Result      string  `json:"result,omitempty"`
	Error       string  `json:"error,omitempty"`
	CreatedAt   string  `json:"created_at"`
	ProcessedAt *string `json:"processed_at,omitempty"`
}

type AgentStatusResponse struct {
	Uptime         string `json:"uptime"`
	ProcessedCount int64  `json:"processed_count"`
	QueueDepth     int    `json:"queue_depth"`
	Model          string `json:"model"`
}

func taskToResponse(t agent.Task) TaskResponse {
	resp := TaskResponse{
		ID:        t.ID,
		Status:    string(t.Status),
		Prompt:    t.Prompt,
		Result:    t.Result,
		Error:     t.Error,
		CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if t.ProcessedAt != nil {
		formatted := t.ProcessedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.ProcessedAt = &formatted
	}

	return resp
}

func statusToResponse(s agent.Status) AgentStatusResponse {
	return AgentStatusResponse{
		Uptime:         s.Uptime.String(),
		ProcessedCount: s.ProcessedCount,
		QueueDepth:     s.QueueDepth,
		Model:          s.Model,
	}
}
