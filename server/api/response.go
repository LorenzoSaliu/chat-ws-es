package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type meta struct {
	Status        int       `json:"status"`
	Message       string    `json:"message"`
	TimeStamp     time.Time `json:"time_stamp"`
	CorrelationId string    `json:"correlation_id"`
}

type response struct {
	Meta   meta        `json:"meta"`
	Result interface{} `json:"result"`
}

func Response(c *fiber.Ctx, err any, status int, message string, correlation_id string) error {
	return c.Status(status).JSON(response{
		Meta: meta{
			Status:        status,
			Message:       message,
			TimeStamp:     time.Now(),
			CorrelationId: correlation_id},
		Result: &fiber.Map{"data": err}})
}
