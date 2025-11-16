package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
	"go.uber.org/zap"
)

// MarketingServiceInterface defines methods for Marketing operations
type MarketingServiceInterface interface {
	GetMarketingData(ctx context.Context, projectID uint) (*services.MarketingData, error)
}

// MarketingHandler handles HTTP requests for marketing data
type MarketingHandler struct {
	marketingService MarketingServiceInterface
}

// NewMarketingHandler creates a new marketing handler
func NewMarketingHandler(marketingService MarketingServiceInterface) *MarketingHandler {
	return &MarketingHandler{
		marketingService: marketingService,
	}
}

// GetMarketing handles GET /api/projects/:id/marketing
func (h *MarketingHandler) GetMarketing(c echo.Context) error {
	ctx := c.Request().Context()

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid project ID")
	}

	marketingData, err := h.marketingService.GetMarketingData(ctx, uint(projectID))
	if err != nil {
		logger.Log.Error("Failed to get marketing data", zap.Error(err), zap.Uint64("project_id", projectID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve marketing data")
	}

	return c.JSON(http.StatusOK, marketingData)
}

