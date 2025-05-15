package auth_middleware

import (
	"github.com/gofiber/fiber/v2"
)

// VerifyMetaWebhook handles the verification of Facebook webhook requests
func MetaVerificationRequestToken(expectedToken string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		mode := c.Query("hub.mode")
		token := c.Query("hub.verify_token")

		if mode != "subscribe" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid mode")
		}
		if token != expectedToken {
			return fiber.NewError(fiber.StatusForbidden, "Invalid token")
		}

		return c.Next()
	}
}
