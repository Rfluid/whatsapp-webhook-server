package webhook_router

import (
	server_model "github.com/Rfluid/whatsapp-webhook-server/src/server/model"
	webhook_service "github.com/Rfluid/whatsapp-webhook-server/src/webhook/service"
)

func Route(server server_model.Config, hook *webhook_service.Config) {
	// Spread middlewares for the Post route
	server.WebhookGroup.Post(
		hook.Path,
		append(hook.PostMiddlewares, hook.Post)..., // Spread middlewares followed by the handler
	)

	// Spread middlewares for the Get route
	server.WebhookGroup.Get(
		hook.Path,
		append(hook.GetMiddlewares, hook.Get)..., // Spread middlewares followed by the handler
	)
}
