package webhook_model

import (
	"slices"

	webhook "github.com/Rfluid/whatsapp-cloud-api/src/webhook"
	"github.com/gofiber/fiber/v2"
)

type HandlerCallback = func(*fiber.Ctx, *webhook.WebhookBody, *webhook.Change) error

// Function that will be executed if one of the contexts matches the real context of the webhook.
type ChangeHandler struct {
	Callback          HandlerCallback
	ExecutionContexts *[]webhook.Field
}

func (h *ChangeHandler) ExecConditionally(ctx *fiber.Ctx, body *webhook.WebhookBody, change *webhook.Change) error {
	if h.ExecutionContexts == nil || slices.Contains(*h.ExecutionContexts, change.Field) {
		if err := h.Callback(ctx, body, change); err != nil {
			return err
		}
	}
	return nil
}
