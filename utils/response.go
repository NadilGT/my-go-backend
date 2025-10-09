package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Operation string `json:"operation"`
	Error     string `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
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

func NewCustomError(c *fiber.Ctx, status int, message string, err error) error {
	resp := fiber.Map{
		"status":  status,
		"message": message,
	}
	if err != nil {
		resp["error"] = err.Error()
	}
	return c.Status(status).JSON(resp)
}
