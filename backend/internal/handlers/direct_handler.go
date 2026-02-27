package handlers

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/suprt/planica_bi/backend/internal/models"
	"github.com/suprt/planica_bi/backend/internal/services"
)

// DirectServiceInterface defines methods for Direct operations
type DirectServiceInterface interface {
	CreateAccount(ctx context.Context, account *models.DirectAccount) error
	GetAccountsByProject(ctx context.Context, projectID uint) ([]*models.DirectAccount, error)
	GetCampaignsWithMetrics(ctx context.Context, projectID uint) ([]services.CampaignWithMetrics, error)
}

// DirectHandler handles HTTP requests for Direct accounts
type DirectHandler struct {
	directService DirectServiceInterface
}

// NewDirectHandler creates a new Direct handler
func NewDirectHandler(directService DirectServiceInterface) *DirectHandler {
	return &DirectHandler{
		directService: directService,
	}
}

// AddDirectAccount handles POST /api/projects/:id/direct-accounts
func (h *DirectHandler) AddDirectAccount(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	var account models.DirectAccount
	if err := c.Bind(&account); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Set project ID from URL
	account.ProjectID = uint(projectID)

	if err := h.directService.CreateAccount(ctx, &account); err != nil {
		// Validation errors should return 400, other errors will be handled by error handler
		if err.Error() == "project_id is required" || err.Error() == "client_login is required" {
			return echo.NewHTTPError(400, err.Error())
		}
		if err.Error() == "account with this ClientLogin already exists for this project" {
			return echo.NewHTTPError(409, err.Error())
		}
		return err
	}

	return c.JSON(201, account)
}

// GetDirectAccounts handles GET /api/projects/:id/direct-accounts
func (h *DirectHandler) GetDirectAccounts(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	accounts, err := h.directService.GetAccountsByProject(ctx, uint(projectID))
	if err != nil {
		return err
	}

	return c.JSON(200, accounts)
}

// GetCampaigns handles GET /api/projects/:id/campaigns
func (h *DirectHandler) GetCampaigns(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid project ID")
	}

	campaigns, err := h.directService.GetCampaignsWithMetrics(ctx, uint(projectID))
	if err != nil {
		return err
	}

	return c.JSON(200, campaigns)
}
