package responses

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type GenericResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}

/// SUCCESS ///

// 200
func OKResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusOK).JSON(GenericResponse{Status: http.StatusOK, Message: "success", Data: data})
}

// 201
func CreatedResponse(c *fiber.Ctx, data *fiber.Map) error {
	return c.Status(http.StatusCreated).JSON(GenericResponse{Status: http.StatusCreated, Message: "success", Data: data})
}

/// ERRORS ///

// 400
func BadRequestResponse(c *fiber.Ctx, data string) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"error": data}})
}

// 401
func UnauthorizedResponse(c *fiber.Ctx, data string) error {
	return c.Status(http.StatusUnauthorized).JSON(GenericResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"error": data}})
}

// 403
func ForbiddenResponse(c *fiber.Ctx, data string) error {
	return c.Status(http.StatusForbidden).JSON(GenericResponse{Status: http.StatusForbidden, Message: "error", Data: &fiber.Map{"error": data}})
}

// 404
func NotFoundResponse(c *fiber.Ctx, data string) error {
	return c.Status(http.StatusNotFound).JSON(GenericResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": data}})
}

// 500
func InternalServerErrorResponse(c *fiber.Ctx, data string) error {
	return c.Status(http.StatusInternalServerError).JSON(GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"error": data}})
}
