package handlers

import (
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/services"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login godoc
// @Summary      User login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body services.LoginRequest true "Login credentials"
// @Success      200  {object} response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}

	tokens, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}
	return response.Success(c, "Login successful", tokens)
}

// Logout godoc
// @Summary      User logout
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	type logoutReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	var req logoutReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if err := h.svc.Logout(c.Context(), req.RefreshToken); err != nil {
		return response.InternalError(c, "Logout failed")
	}
	return response.Success(c, "Logged out successfully", nil)
}

// Refresh godoc
// @Summary      Refresh access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	type refreshReq struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	var req refreshReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	tokens, err := h.svc.RefreshTokens(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}
	return response.Success(c, "Token refreshed", tokens)
}

// Me godoc
// @Summary      Get current user profile
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims, ok := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	if !ok {
		return response.Unauthorized(c, "Unauthorized")
	}
	return response.Success(c, "Profile retrieved", fiber.Map{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}
