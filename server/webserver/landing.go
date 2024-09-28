package webserver

import "github.com/gofiber/fiber/v2"

func (ws *WebServer) ServeLandingPage(ctx *fiber.Ctx) error {
	return ctx.Render("landing", fiber.Map{})
}
