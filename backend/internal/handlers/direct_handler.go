package handlers

import (
	"github.com/labstack/echo/v4"
)

// DirectHandler handles HTTP requests for Direct accounts
type DirectHandler struct {
	// TODO: add service
}

// NewDirectHandler creates a new Direct handler
func NewDirectHandler() *DirectHandler {
	return &DirectHandler{}
}

// AddDirectAccount handles POST /api/projects/:id/direct-accounts
func (h *DirectHandler) AddDirectAccount(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}

// GetDirectAccounts handles GET /api/projects/:id/direct-accounts
func (h *DirectHandler) GetDirectAccounts(c echo.Context) error {
	// TODO: implement
	return c.NoContent(501) // Not Implemented
}
