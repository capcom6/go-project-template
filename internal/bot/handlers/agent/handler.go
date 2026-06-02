package agent

import (
	"fmt"
	"strings"

	"github.com/capcom6/go-project-template/internal/agent"
	"github.com/capcom6/go-project-template/internal/bot/handler"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
)

type Handler struct {
	agentSvc *agent.Service
	logger   *zap.Logger
}

func New(agentSvc *agent.Service, logger *zap.Logger) handler.Handler {
	return &Handler{
		agentSvc: agentSvc,
		logger:   logger,
	}
}

func (h *Handler) Register(router *telegofx.Router) {
	router.Handle(
		h.handleAgent,
		th.CommandEqual("agent"),
		th.AnyMessageWithFrom(),
	)
	router.Handle(
		h.handleAgentTask,
		th.CommandEqual("agent_task"),
		th.AnyMessageWithFrom(),
	)
}

func (h *Handler) handleAgent(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	status := h.agentSvc.Status()
	text := fmt.Sprintf(
		"AI Agent Status:\n"+
			"Model: %s\n"+
			"Uptime: %s\n"+
			"Processed: %d\n"+
			"Queue depth: %d",
		status.Model,
		status.Uptime,
		status.ProcessedCount,
		status.QueueDepth,
	)

	return h.reply(ctx, update.Message.Chat.ID, text)
}

func (h *Handler) handleAgentTask(ctx *th.Context, update telego.Update) error {
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	prompt := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/agent_task"))
	if prompt == "" {
		return h.reply(ctx, update.Message.Chat.ID, "Usage: /agent_task <your prompt>")
	}

	task, err := h.agentSvc.Enqueue(ctx, prompt)
	if err != nil {
		h.logger.Error("failed to enqueue task", zap.Error(err))
		return h.reply(ctx, update.Message.Chat.ID, "Failed to enqueue task. Try again later.")
	}

	return h.reply(ctx, update.Message.Chat.ID,
		fmt.Sprintf("Task #%d enqueued: %s", task.ID, task.Prompt))
}

func (h *Handler) reply(ctx *th.Context, chatID int64, text string) error {
	_, err := ctx.Bot().SendMessage(ctx, &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: chatID},
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("send telegram message: %w", err)
	}

	return nil
}
