package server_model

import "github.com/gofiber/fiber/v2"

type Config struct {
	Path         string
	App          *fiber.App
	WebhookGroup fiber.Router
}

func (c *Config) SetGroup() {
	c.WebhookGroup = c.App.Group(c.Path)
}
