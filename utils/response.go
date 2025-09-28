package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Operation string `json:"operation"`
	Error     string `json:"error,omitempty"`
}

func SendErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	response := Response{
		Operation: "Failed",
		Error:     errorMessage,
	}
	return c.Status(statusCode).JSON(response)
}

func SendSuccessResponse(c *fiber.Ctx) error {
	response := Response{
		Operation: "Success",
		Error:     "",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
