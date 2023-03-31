package routes

import (
	Controller "secusend/controllers"

	"github.com/gofiber/fiber/v2"
)

func NoteRouter(router fiber.Router) {
	router.Post("/postNote", Controller.PostNote())
	router.Get("/getNote", Controller.GetNote())
}
