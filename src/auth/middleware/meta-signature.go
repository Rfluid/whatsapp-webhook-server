package auth_middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// VerifyMetaSignature is a Fiber middleware that verifies the X-Hub-Signature-256 header
func VerifyMetaSignature(appSecret string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		signatureHeader := c.Get("X-Hub-Signature-256")
		if signatureHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing X-Hub-Signature-256 header")
		}

		const prefix = "sha256="
		if !strings.HasPrefix(signatureHeader, prefix) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature format")
		}

		signature := signatureHeader[len(prefix):]
		body := c.Body()

		mac := hmac.New(sha256.New, []byte(appSecret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)
		expectedSignature := hex.EncodeToString(expectedMAC)

		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid signature")
		}

		return c.Next()
	}
}
