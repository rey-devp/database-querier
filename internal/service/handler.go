package service

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"database-querier-agent/internal/agent"
	"database-querier-agent/internal/logger"
	"database-querier-agent/internal/memory"
)

type Handler struct {
	agent *agent.Agent
	store *memory.Store
}

func NewHandler(a *agent.Agent, s *memory.Store) *Handler {
	return &Handler{
		agent: a,
		store: s,
	}
}

func (h *Handler) HandleQuery(c *fiber.Ctx) error {
	logger.Info("REQUEST", "Received POST /query", "ip", c.IP())

	var task memory.Task
	if err := c.BodyParser(&task); err != nil {
		logger.Warn("REQUEST", "Invalid request body", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Save task to mock store
	h.store.SaveTask(&task)

	// Process task
	res, err := h.agent.ProcessTask(context.Background(), task.ID)
	
	if err != nil {
		// Output the ErrorResponse that was saved to store
		errRes, _ := h.store.GetResult(task.ID)
		return c.Status(fiber.StatusBadRequest).JSON(errRes)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) HandleHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
