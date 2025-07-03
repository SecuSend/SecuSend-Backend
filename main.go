package main

import (
	"flag"
	"log"
	"secusend/configs"
	"secusend/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New()

	//DB
	configs.ConnectDB()

	// Middleware
	app.Use(limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	// app.Use(csrf.New())
	//todo add compress?

	// Routes
	app.Get("/", HealthCheck)
	routes.NoteRouter(app.Group("/api/note"))

	log.Fatal(app.Listen(":3000"))
}

func HealthCheck(c *fiber.Ctx) error {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	if err := c.JSON(res); err != nil {
		return err
	}

	return nil
}
