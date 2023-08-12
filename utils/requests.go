package utils

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

// ParseRequestBody parses the body of a request.
func ParseRequestBody(ctx *fiber.Ctx, body any) {
	err := ctx.BodyParser(&body)
	if err != nil {
		err = ctx.Status(500).SendString(err.Error())
		if err != nil {
			log.Fatalln("Error sending response when parsing request body:", err)
		}
	}
}
