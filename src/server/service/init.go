package server_service

import (
	server_model "github.com/Rfluid/whatsapp-webhook-server/src/server/model"
	webhook_router "github.com/Rfluid/whatsapp-webhook-server/src/webhook/router"
	webhook_service "github.com/Rfluid/whatsapp-webhook-server/src/webhook/service"
)

func Bootstrap(server *server_model.Config, hook *webhook_service.Config) {
	server.SetGroup()
	webhook_router.Route(*server, hook)
}
