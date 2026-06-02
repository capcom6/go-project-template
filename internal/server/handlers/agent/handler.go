package agent

import (
	"github.com/capcom6/go-project-template/internal/agent"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	handler.Base

	agentSvc *agent.Service
}

func New(service *agent.Service, validator *validator.Validate) handler.Handler {
	return &Handler{
		Base: handler.Base{
			Validator: validator,
		},
		agentSvc: service,
	}
}

func (h *Handler) Register(r fiber.Router) {
	g := r.Group("/agent")

	g.Get("/status", h.getStatus)
	g.Get("/tasks", h.getTasks)
	g.Post("/tasks", validation.DecorateWithBodyEx(h.Validator, h.enqueue))
}

// @Summary		Get agent status
// @Description	Returns the current status of the AI agent
// @Tags			agent
// @Produce		json
// @Success		200	{object}	AgentStatusResponse
// @Failure		500	{object}	fiberfx.ErrorResponse
// @Router			/agent/status [get]
func (h *Handler) getStatus(c *fiber.Ctx) error {
	status := h.agentSvc.Status()
	return c.JSON(statusToResponse(status))
}

// @Summary		List tasks
// @Description	Returns all tasks in the agent queue
// @Tags			agent
// @Produce		json
// @Success		200	{object}	[]TaskResponse
// @Failure		500	{object}	fiberfx.ErrorResponse
// @Router			/agent/tasks [get]
func (h *Handler) getTasks(c *fiber.Ctx) error {
	tasks, err := h.agentSvc.Tasks()
	if err != nil {
		return err
	}

	responses := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		responses[i] = taskToResponse(t)
	}

	return c.JSON(responses)
}

// @Summary		Enqueue task
// @Description	Enqueues a new prompt for the AI agent to process
// @Tags			agent
// @Accept			json
// @Produce		json
// @Param			request	body		EnqueueRequest	true	"Prompt to process"
// @Success		200		{object}	TaskResponse
// @Failure		400		{object}	fiberfx.ErrorResponse
// @Failure		500		{object}	fiberfx.ErrorResponse
// @Router			/agent/tasks [post]
func (h *Handler) enqueue(c *fiber.Ctx, req *EnqueueRequest) error {
	task, err := h.agentSvc.Enqueue(c.Context(), req.Prompt)
	if err != nil {
		return err
	}

	return c.JSON(taskToResponse(task))
}
