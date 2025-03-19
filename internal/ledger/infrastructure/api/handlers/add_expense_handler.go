package handlers

import (
	"lenavire/internal/ledger/application/commands/add_expense"
	"lenavire/internal/ledger/infrastructure/api/dto"

	"github.com/gofiber/fiber/v2"
)

type AddExpenseHandler struct {
	CommandHandler *commands.AddExpenseCommandHandler
}

func NewAddExpenseHandler(commandHandler *commands.AddExpenseCommandHandler) *AddExpenseHandler {
	return &AddExpenseHandler{CommandHandler: commandHandler}
}

func (h *AddExpenseHandler) AddExpense(c *fiber.Ctx) error {
	var req dto.AddExpenseRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.OperationType == "income" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "income operation type is not allowed"})
	}

	command := commands.NewAddExpenseCommand(req.Amount)
	err := h.CommandHandler.Execute(command)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Expense added successfully"})
}
