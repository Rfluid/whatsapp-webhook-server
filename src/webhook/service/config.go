package webhook_service

import (
	"sync"

	wh_model "github.com/Rfluid/whatsapp-cloud-api/src/webhook/model"
	webhook_model "github.com/Rfluid/whatsapp-webhook-server/src/webhook/model"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Path            string
	ChangeHandlers  []webhook_model.ChangeHandler
	CtxHandler      func(ctx *fiber.Ctx, body *wh_model.WebhookBody) error
	PostMiddlewares [](func(ctx *fiber.Ctx) error)
	GetMiddlewares  [](func(ctx *fiber.Ctx) error)
}

// Executes conccurently the change handlers for each change in each entry in the body. The entries are executed concurrently as are the changes and the change handlers.
func (c *Config) ExecConditionally(ctx *fiber.Ctx, body *wh_model.WebhookBody) error {
	var err error = nil
	var entryWg sync.WaitGroup
	entryErrCh := make(chan error, len(body.Entry))

	for _, entry := range body.Entry {
		entryWg.Add(1)
		go func(entry wh_model.Entry) {
			defer entryWg.Done()

			var entryErr error = nil
			var changeWg sync.WaitGroup
			changeCh := make(chan error, len(entry.Changes))

			for _, change := range entry.Changes {
				changeWg.Add(1)
				go func(change wh_model.Change) {
					defer changeWg.Done()

					var changeHandlerErr error = nil
					var changeHandlerWg sync.WaitGroup
					changeHandlerCh := make(chan error, len(c.ChangeHandlers))

					for _, h := range c.ChangeHandlers {
						changeHandlerWg.Add(1)
						go func() {
							defer changeHandlerWg.Done()
							changeHandlerCh <- h.ExecConditionally(ctx, body, &change)
						}()
					}

					go func() {
						changeHandlerWg.Wait()
						close(changeHandlerCh)
					}()

					for errInCh := range changeHandlerCh {
						if errInCh != nil {
							changeHandlerErr = errInCh
						}
					}

					changeCh <- changeHandlerErr
				}(change)
			}

			go func() {
				changeWg.Wait()
				close(changeCh)
			}()

			for errInCh := range changeCh {
				if errInCh != nil {
					entryErr = errInCh
				}
			}

			entryErrCh <- entryErr
		}(entry)
	}

	go func() {
		entryWg.Wait()
		close(entryErrCh)
	}()

	for errInCh := range entryErrCh {
		if errInCh != nil {
			err = errInCh
		}
	}

	return err
}

func (c *Config) Exec(ctx *fiber.Ctx, body *wh_model.WebhookBody) error {
	err := c.CtxHandler(ctx, body)
	if err != nil {
		return err
	}
	return c.ExecConditionally(ctx, body)
}

// @Summary		Handles Webhooks.
// @Description	Executes the context handler then the change handlers. If any error is thrown this function will also throw an error.
// @Tags			Webhook
// @Accept			json
// @Produce		json
// @Param			input	body	wh_model.WebhookBody	true	"Content sent by WhatsApp Cloud API."
// @Success		200 "Valid webhook endpoint."
// @Router			/{webhook_path} [post] // Change the docs adding the real path here.
func (c *Config) Post(ctx *fiber.Ctx) error {
	var body wh_model.WebhookBody
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	err := c.Exec(ctx, &body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// @Summary		Verify Webhook.
// @Description	Used by meta to verify if it is a valid webhook endpoint.
// @Tags			Webhook
// @Accept			json
// @Produce		json
// @Param			hub.mode			query	string	true	"Subscription mode, always set to 'subscribe'"
// @Param			hub.challenge		query	int		true	"A challenge integer that must be returned to confirm the webhook"
// @Param			hub.verify_token	query	string	true	"A string used for validation, defined in the Webhooks setup in the App Dashboard"
// @Success		200 {string} string	"hub.challenge returned as a string."
// @Router			/{webhook_path} [get] // Change the docs adding the real path here.
func (c *Config) Get(ctx *fiber.Ctx) error {
	hubChallenge := ctx.Query("hub.challenge")
	return ctx.Status(fiber.StatusOK).SendString(hubChallenge)
}
