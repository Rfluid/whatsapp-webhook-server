package server_service

import (
	server_model "github.com/Rfluid/whatsapp-webhook-server/src/server/model"
	"github.com/gofiber/fiber/v2"
)

func DefaultConfig(app *fiber.App) *server_model.Config {
	return NewConfig(app, "/webhook")
}

func NewConfig(app *fiber.App, path string) *server_model.Config {
	config := server_model.Config{
		Path: path,
		App:  app,
	}

	config.SetGroup()

	return &config
}
