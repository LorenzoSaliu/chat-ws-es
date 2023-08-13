package main

import (
	"log"
	"net/http"
	"os"

	"server/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	logfile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	app.Use(logger.New(logger.Config{Output: logfile}))

	app.Get("/chat/health-check", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "server is up and running"})
	})

	v1 := app.Group("/chat/v1")

	v1.Post("/signin", api.SignInHandler)
	v1.Post("/login", api.LogInHandler)
	v1.Get("/user/:user_id", api.GetUserHandler)

	//start server localhost:3000
	log.Fatal(app.Listen(":3000"))

}
