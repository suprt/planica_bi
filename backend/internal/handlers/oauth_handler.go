package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/config"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/logger"
	"go.uber.org/zap"
)

const (
	yandexOAuthAuthorizeURL = "https://oauth.yandex.ru/authorize"
	yandexOAuthTokenURL     = "https://oauth.yandex.ru/token"
)

// TokenResponse represents the response from Yandex OAuth token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// TokenErrorResponse represents an error response from Yandex OAuth token endpoint
type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// OAuthHandler handles OAuth authentication with Yandex
type OAuthHandler struct {
	cfg *config.Config
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(cfg *config.Config) *OAuthHandler {
	return &OAuthHandler{
		cfg: cfg,
	}
}

// InitiateAuth handles GET /api/oauth/yandex
// Redirects user to Yandex OAuth authorization page
// Documentation: https://yandex.ru/dev/id/doc/ru/concepts/ya-oauth-intro
func (h *OAuthHandler) InitiateAuth(c echo.Context) error {
	if h.cfg.YandexClientID == "" {
		return echo.NewHTTPError(500, "Yandex Client ID is not configured")
	}

	// Build redirect URI
	redirectURI := fmt.Sprintf("%s/api/oauth/yandex/callback", h.cfg.AppURL)

	// Build authorization URL with required parameters
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", h.cfg.YandexClientID)
	params.Set("redirect_uri", redirectURI)

	// Set OAuth scopes from config (space-separated list)
	// If scopes are not specified, Yandex will use scopes configured during app registration
	// IMPORTANT: Scopes must match those specified during app registration in Yandex OAuth
	// Check your app settings at: https://oauth.yandex.ru/client/<client_id>/info
	if h.cfg.YandexOAuthScopes != "" {
		params.Set("scope", h.cfg.YandexOAuthScopes)
	}

	authURL := fmt.Sprintf("%s?%s", yandexOAuthAuthorizeURL, params.Encode())

	// Log the authorization URL for debugging
	if logger.Log != nil {
		logger.Log.Info("OAuth authorization initiated",
			zap.String("url", authURL),
			zap.String("redirect_uri", redirectURI),
		)
	}

	return c.Redirect(http.StatusFound, authURL)
}

// HandleCallback handles GET /api/oauth/yandex/callback
// Processes OAuth callback from Yandex, exchanges code for token, and redirects to frontend
func (h *OAuthHandler) HandleCallback(c echo.Context) error {
	ctx := c.Request().Context()

	if logger.Log != nil {
		logger.Log.Info("OAuth callback received",
			zap.String("url", c.Request().URL.String()),
			zap.String("query", c.Request().URL.RawQuery),
		)
	}

	// Get authorization code from query parameter
	code := c.QueryParam("code")
	if code == "" {
		if logger.Log != nil {
			logger.Log.Warn("OAuth callback: authorization code is missing")
		}
		return echo.NewHTTPError(400, "Authorization code is missing")
	}

	// Check for error from Yandex
	if errorParam := c.QueryParam("error"); errorParam != "" {
		errorDescription := c.QueryParam("error_description")
		if logger.Log != nil {
			logger.Log.Error("OAuth callback error from Yandex",
				zap.String("error", errorParam),
				zap.String("description", errorDescription),
			)
		}
		// Redirect to frontend with error parameter
		frontendURL := h.cfg.FrontendURL
		if frontendURL == "" {
			frontendURL = h.cfg.AppURL
		}
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard?oauth=error&error=%s", frontendURL, url.QueryEscape(errorParam)))
	}

	if logger.Log != nil {
		logger.Log.Info("OAuth callback: code received, exchanging for token")
	}

	// Exchange code for access token
	tokenResp, err := h.exchangeCodeForToken(ctx, code)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("OAuth callback: failed to exchange code for token",
				zap.Error(err),
			)
		}
		// Redirect to frontend with error parameter
		frontendURL := h.cfg.FrontendURL
		if frontendURL == "" {
			frontendURL = h.cfg.AppURL
		}
		return c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard?oauth=error", frontendURL))
	}

	// Save token to .env file (for MVP)
	// TODO: In production, save to database instead
	if err := h.saveTokenToEnv(tokenResp.AccessToken); err != nil {
		if logger.Log != nil {
			logger.Log.Warn("Failed to save OAuth token to .env file",
				zap.Error(err),
			)
		}
		// Don't fail the request, but log the warning
	} else {
		if logger.Log != nil {
			logger.Log.Info("OAuth token saved to .env file")
		}
	}

	// Log token for debugging
	if logger.Log != nil {
		tokenPreview := ""
		if len(tokenResp.AccessToken) > 20 {
			tokenPreview = tokenResp.AccessToken[:20] + "..."
		} else {
			tokenPreview = tokenResp.AccessToken
		}
		logger.Log.Info("OAuth token received",
			zap.String("token_preview", tokenPreview),
			zap.Int("expires_in", tokenResp.ExpiresIn),
			zap.Bool("has_refresh_token", tokenResp.RefreshToken != ""),
		)
	}

	// Redirect to frontend with success parameter
	frontendURL := h.cfg.FrontendURL
	if frontendURL == "" {
		// Fallback to AppURL if FrontendURL is not configured
		frontendURL = h.cfg.AppURL
	}
	if logger.Log != nil {
		logger.Log.Info("OAuth callback: redirecting to frontend",
			zap.String("frontend_url", frontendURL),
		)
	}
	return c.Redirect(http.StatusFound, fmt.Sprintf("%s/dashboard?oauth=success", frontendURL))
}

// GetOAuthStatus handles GET /api/oauth/status
// Returns OAuth authorization status by validating the token with Yandex API
func (h *OAuthHandler) GetOAuthStatus(c echo.Context) error {
	ctx := c.Request().Context()

	// If no token in config, definitely not authorized
	if h.cfg.YandexOAuthToken == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"authorized": false,
			"has_token":  false,
		})
	}

	// Validate token by making a request to Yandex API
	// Using a simple info endpoint to check if token is valid
	isValid, err := h.validateToken(ctx, h.cfg.YandexOAuthToken)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Warn("Failed to validate OAuth token",
				zap.Error(err),
			)
		}
		// If validation fails, assume token is invalid
		isValid = false
	}

	if logger.Log != nil {
		tokenPreview := ""
		if len(h.cfg.YandexOAuthToken) > 20 {
			tokenPreview = h.cfg.YandexOAuthToken[:20] + "..."
		} else {
			tokenPreview = h.cfg.YandexOAuthToken
		}
		logger.Log.Info("OAuth status check",
			zap.String("token_preview", tokenPreview),
			zap.Bool("is_valid", isValid),
			zap.Bool("has_token", h.cfg.YandexOAuthToken != ""),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"authorized": isValid,
		"has_token":  h.cfg.YandexOAuthToken != "",
	})
}

// validateToken validates OAuth token by making a request to Yandex ID API
// Returns true if token is valid, false otherwise
func (h *OAuthHandler) validateToken(ctx context.Context, token string) (bool, error) {
	// Use Yandex ID API info endpoint to validate token
	// Documentation: https://yandex.ru/dev/id/doc/ru/user-information
	infoURL := "https://login.yandex.ru/info"

	req, err := http.NewRequestWithContext(ctx, "GET", infoURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{
		Timeout: 5 * time.Second, // Short timeout for status check
	}

	resp, err := client.Do(req)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Debug("Token validation request failed",
				zap.Error(err),
			)
		}
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read body to check for error details
	body, _ := io.ReadAll(resp.Body)

	if logger.Log != nil {
		logger.Log.Debug("Token validation response",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_preview", string(body)[:min(len(string(body)), 100)]),
		)
	}

	// If status is 200, token is valid
	// If status is 401, token is invalid or expired
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// exchangeCodeForToken exchanges authorization code for access token
// Documentation: https://yandex.ru/dev/id/doc/ru/access
func (h *OAuthHandler) exchangeCodeForToken(ctx context.Context, code string) (*TokenResponse, error) {
	if h.cfg.YandexClientID == "" || h.cfg.YandexClientSecret == "" {
		return nil, fmt.Errorf("yandex OAuth credentials are not configured")
	}

	// Build redirect URI (must match the one used in InitiateAuth)
	redirectURI := fmt.Sprintf("%s/api/oauth/yandex/callback", h.cfg.AppURL)

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", h.cfg.YandexClientID)
	data.Set("client_secret", h.cfg.YandexClientSecret)
	data.Set("redirect_uri", redirectURI)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", yandexOAuthTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Log request for debugging (without sensitive data)
	if logger.Log != nil {
		logger.Log.Debug("OAuth token exchange request",
			zap.String("url", yandexOAuthTokenURL),
			zap.String("redirect_uri", redirectURI),
			zap.String("client_id", h.cfg.YandexClientID),
		)
	}

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("OAuth token exchange: request failed",
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Error("OAuth token exchange: failed to read response",
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		// Try to parse error response as JSON
		var errorResp TokenErrorResponse
		if jsonErr := json.Unmarshal(body, &errorResp); jsonErr == nil && errorResp.Error != "" {
			if logger.Log != nil {
				logger.Log.Error("OAuth token exchange failed",
					zap.Int("status", resp.StatusCode),
					zap.String("error", errorResp.Error),
					zap.String("error_description", errorResp.ErrorDescription),
				)
			}
			return nil, fmt.Errorf("OAuth token exchange failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
		}

		// Fallback to raw body if JSON parsing failed
		if logger.Log != nil {
			logger.Log.Error("OAuth token exchange failed",
				zap.Int("status", resp.StatusCode),
				zap.String("response", string(body)),
			)
		}
		return nil, fmt.Errorf("OAuth token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("access token is missing in response")
	}

	return &tokenResp, nil
}

// saveTokenToEnv saves OAuth token to .env file
// For MVP: stores token in .env file (can be moved to DB later)
func (h *OAuthHandler) saveTokenToEnv(accessToken string) error {
	return config.UpdateEnvFile("YANDEX_OAUTH_TOKEN", accessToken)
}
