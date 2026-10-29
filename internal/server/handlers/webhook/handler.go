package webhook

import (
	"context"

	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/telegofx"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Handler serves the Telegram webhook endpoint. The route is mounted in
// every mode; only webhook mode has a handler ready to accept the payload,
// polling mode answers 500 (telegofx ErrHandlerNotReady) because Telegram
// never sends updates to this service.
type Handler struct {
	feeder telegofx.WebhookHandler

	logger *zap.Logger
}

// New creates a new webhook HTTP handler.
func New(feeder telegofx.WebhookHandler, logger *zap.Logger) handler.Handler {
	return &Handler{
		feeder: feeder,
		logger: logger,
	}
}

// Register registers the webhook HTTP route.
func (h *Handler) Register(r fiber.Router) {
	r.Post("/webhook", h.handle)
}

//	@Summary		Handle Telegram update
//	@Description	Accepts a raw Telegram webhook update and pushes it into the ingestion pipeline. Polled updates are not sent here, so the route answers 500 while the bot runs in polling mode. The route performs no request authentication.
//	@Tags			webhook
//	@Accept			json
//	@Produce		json
//	@Param			update	body		object	true	"Telegram update"
//	@Success		200		{object}	Response
//	@Failure		500		{object}	ErrorResponse
//	@Router			/webhook [post]
//
// handle decodes the Telegram update and pushes it into the ingestion
// pipeline. The update is processed after the request returns, so the request
// context is detached (telego webhook guidance). Every feeder failure answers
// 500, including telegofx ErrHandlerNotReady while the bot runs in polling
// mode or after it has stopped.
func (h *Handler) handle(c *fiber.Ctx) error {
	if err := h.feeder(context.WithoutCancel(c.Context()), c.Body()); err != nil {
		h.logger.Warn("webhook update decode failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: "failed to decode update"})
	}

	return c.JSON(Response{OK: true})
}
