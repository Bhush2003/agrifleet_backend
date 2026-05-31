package middleware

import (
	"strings"

	"github.com/agrifleet/backend/internal/models"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

const LocalsUserKey = "user_claims"

// AuthMiddleware validates the Bearer JWT token on protected routes.
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Authorization header is required")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Unauthorized(c, "Invalid authorization format. Use: Bearer <token>")
		}

		claims, err := pkgjwt.ParseToken(parts[1], jwtSecret)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		c.Locals(LocalsUserKey, claims)
		return c.Next()
	}
}

// RequireRoles restricts access to users with one of the specified roles.
func RequireRoles(roles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(LocalsUserKey).(*pkgjwt.Claims)
		if !ok || claims == nil {
			return response.Unauthorized(c, "Unauthorized")
		}

		for _, role := range roles {
			if models.UserRole(claims.Role) == role {
				return c.Next()
			}
		}
		return response.Forbidden(c, "You do not have permission to access this resource")
	}
}
