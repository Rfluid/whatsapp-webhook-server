package webhook_model

import (
	"slices"

	wh_model "github.com/Rfluid/whatsapp-cloud-api/src/webhook/model"
	"github.com/gofiber/fiber/v2"
)

type HandlerCallback = func(*fiber.Ctx, *wh_model.WebhookBody, *wh_model.Change) error

// Function that will be executed if one of the contexts matches the real context of the webhook.
type ChangeHandler struct {
	Callback          HandlerCallback
	ExecutionContexts *[]wh_model.Field
}

func (h *ChangeHandler) ExecConditionally(ctx *fiber.Ctx, body *wh_model.WebhookBody, change *wh_model.Change) error {
	if h.ExecutionContexts == nil || slices.Contains(*h.ExecutionContexts, change.Field) {
		if err := h.Callback(ctx, body, change); err != nil {
			return err
		}
	}
	return nil
}
