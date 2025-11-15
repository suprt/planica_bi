package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	yandexDirectAPIURL     = "https://api.direct.yandex.com/json/v5"
	yandexDirectSandboxURL = "https://api-sandbox.direct.yandex.com/json/v5"
)

// YandexDirectClient handles integration with Yandex.Direct API
type YandexDirectClient struct {
	token       string
	clientLogin string
	httpClient  *http.Client
	useSandbox  bool
	baseURL     string // For testing: allows overriding base URL
}

// NewYandexDirectClient creates a new Direct client
// useSandbox: if true, uses sandbox environment (https://api-sandbox.direct.yandex.com)
func NewYandexDirectClient(token, clientLogin string, useSandbox bool) *YandexDirectClient {
	baseURL := yandexDirectAPIURL
	if useSandbox {
		baseURL = yandexDirectSandboxURL
	}
	return &YandexDirectClient{
		token:       token,
		clientLogin: clientLogin,
		useSandbox:  useSandbox,
		baseURL:     baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewYandexDirectClientWithURL creates a new Direct client with custom base URL (for testing)
func NewYandexDirectClientWithURL(token, clientLogin, baseURL string) *YandexDirectClient {
	return &YandexDirectClient{
		token:       token,
		clientLogin: clientLogin,
		baseURL:     baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Campaign represents a Yandex Direct campaign
type Campaign struct {
	Id   int64  `json:"Id"`
	Name string `json:"Name"`
}

// CampaignsResponse represents response from campaigns.get
type CampaignsResponse struct {
	Result struct {
		Campaigns []Campaign `json:"Campaigns"`
	} `json:"result"`
	Error *APIError `json:"error"`
}

// ReportResponse represents response from reports
type ReportResponse struct {
	Result struct {
		Report struct {
			Rows []ReportRow `json:"Rows"`
		} `json:"Report"`
	} `json:"result"`
	Error *APIError `json:"error"`
}

// ReportRow represents a row in report data
type ReportRow struct {
	CampaignId   int64   `json:"CampaignId"`
	CampaignName string  `json:"CampaignName"`
	Impressions  int64   `json:"Impressions"`
	Clicks       int64   `json:"Clicks"`
	Cost         float64 `json:"Cost"`
	CTR          float64 `json:"Ctr"`
	AvgCpc       float64 `json:"AvgCpc"`
	Conversions  int64   `json:"Conversions,omitempty"`
	CPA          float64 `json:"Cpa,omitempty"`
}

// APIError represents an error from Yandex Direct API
type APIError struct {
	ErrorString string `json:"error_string"`
	ErrorCode   int    `json:"error_code"`
}

// GetCampaigns retrieves campaigns for an account
// Documentation: https://yandex.ru/dev/direct/doc/ref-v5/campaigns/get.html
func (c *YandexDirectClient) GetCampaigns(ctx context.Context) ([]Campaign, error) {
	baseURL := c.baseURL
	if baseURL == "" {
		// Fallback for old clients without baseURL set
		baseURL = yandexDirectAPIURL
		if c.useSandbox {
			baseURL = yandexDirectSandboxURL
		}
	}
	url := baseURL + "/campaigns"

	requestBody := map[string]interface{}{
		"method": "get",
		"params": map[string]interface{}{
			"SelectionCriteria": map[string]interface{}{},
			"FieldNames":        []string{"Id", "Name"},
		},
	}

	var response CampaignsResponse
	if err := c.makeRequest(ctx, url, requestBody, &response); err != nil {
		return nil, fmt.Errorf("failed to get campaigns: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", response.Error.ErrorString, response.Error.ErrorCode)
	}

	return response.Result.Campaigns, nil
}

// GetCampaignReport retrieves report data for campaigns
// Documentation: https://yandex.ru/dev/direct/doc/reports/reports.html
func (c *YandexDirectClient) GetCampaignReport(ctx context.Context, dateFrom, dateTo string) ([]ReportRow, error) {
	baseURL := c.baseURL
	if baseURL == "" {
		// Fallback for old clients without baseURL set
		baseURL = yandexDirectAPIURL
		if c.useSandbox {
			baseURL = yandexDirectSandboxURL
		}
	}
	url := baseURL + "/reports"

	requestBody := map[string]interface{}{
		"params": map[string]interface{}{
			"SelectionCriteria": map[string]interface{}{},
			"FieldNames": []string{
				"CampaignId",
				"CampaignName",
				"Impressions",
				"Clicks",
				"Cost",
				"Ctr",
				"AvgCpc",
				"Conversions",
				"Cpa",
			},
			"ReportName":      "Campaign Report",
			"ReportType":      "CAMPAIGN_PERFORMANCE_REPORT",
			"DateRangeType":   "CUSTOM_DATE",
			"Format":          "TSV",
			"IncludeVAT":      "NO",
			"IncludeDiscount": "NO",
		},
	}

	// Add date range
	if dateFrom != "" && dateTo != "" {
		params := requestBody["params"].(map[string]interface{})
		params["DateFrom"] = dateFrom
		params["DateTo"] = dateTo
	}

	var response ReportResponse
	if err := c.makeRequest(ctx, url, requestBody, &response); err != nil {
		return nil, fmt.Errorf("failed to get campaign report: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("API error: %s (code: %d)", response.Error.ErrorString, response.Error.ErrorCode)
	}

	return response.Result.Report.Rows, nil
}

// makeRequest performs HTTP request to Yandex Direct API
func (c *YandexDirectClient) makeRequest(ctx context.Context, url string, requestBody map[string]interface{}, response interface{}) error {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers according to documentation
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Client-Login", c.clientLogin)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept-Language", "ru")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
