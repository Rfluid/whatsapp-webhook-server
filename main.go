package main

import (
	server_model "github.com/Rfluid/whatsapp-webhook-server/src/server/model"
	server_service "github.com/Rfluid/whatsapp-webhook-server/src/server/service"
	webhook_service "github.com/Rfluid/whatsapp-webhook-server/src/webhook/service"
	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App, server *server_model.Config, hook *webhook_service.Config) {
	server_service.Bootstrap(server, hook)
}
