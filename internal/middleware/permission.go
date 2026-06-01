package middleware

import (
	"github.com/agrifleet/backend/internal/models"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// RequirePermission checks if the authenticated user has a specific permission.
func RequirePermission(perm models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(LocalsUserKey).(*pkgjwt.Claims)
		if !ok || claims == nil {
			return response.Unauthorized(c, "Unauthorized")
		}
		role := models.UserRole(claims.Role)
		if !models.HasPermission(role, perm) {
			return response.Forbidden(c, "You do not have permission: "+string(perm))
		}
		return c.Next()
	}
}

// CheckMonthLock is a middleware that can be injected to block edits on locked months.
// The actual lock check is done in the service layer using the MonthLock model.
func CheckMonthLock() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Month lock enforcement is handled at service layer
		// This middleware can be extended to do a quick DB check
		return c.Next()
	}
}
