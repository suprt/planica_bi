package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	yandexMetricaAPIURL        = "https://api-metrika.yandex.net/stat/v1/data"
	yandexMetricaManagementURL = "https://api-metrika.yandex.net/management/v1"
)

// YandexMetricaClient handles integration with Yandex.Metrica API
type YandexMetricaClient struct {
	token      string
	httpClient *http.Client
	baseURL    string // For testing: allows overriding base URL
}

// NewYandexMetricaClient creates a new Metrica client
// token: OAuth token for authentication
func NewYandexMetricaClient(token string) *YandexMetricaClient {
	return &YandexMetricaClient{
		token:   token,
		baseURL: yandexMetricaAPIURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewYandexMetricaClientWithURL creates a new Metrica client with custom base URL (for testing)
func NewYandexMetricaClientWithURL(token, baseURL string) *YandexMetricaClient {
	return &YandexMetricaClient{
		token:   token,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// MetricsResponse represents response from Metrica API
type MetricsResponse struct {
	Data []MetricsData `json:"data"`
}

// MetricsData represents metrics data row
type MetricsData struct {
	Dimensions []Dimension `json:"dimensions"`
	Metrics    []float64   `json:"metrics"`
}

// Dimension represents a dimension value
type Dimension struct {
	Name string `json:"name"`
}

// AgeMetricsResponse represents response for age breakdown
type AgeMetricsResponse struct {
	Data []AgeMetricsData `json:"data"`
}

// AgeMetricsData represents age metrics data row
type AgeMetricsData struct {
	Dimensions []Dimension `json:"dimensions"`
	Metrics    []float64   `json:"metrics"`
}

// ConversionsResponse represents response for conversions
type ConversionsResponse struct {
	Data []ConversionsData `json:"data"`
}

// ConversionsData represents conversions data row
type ConversionsData struct {
	Dimensions []Dimension `json:"dimensions"`
	Metrics    []float64   `json:"metrics"`
}

// MetricsResult represents parsed metrics result
type MetricsResult struct {
	Visits              int64   `json:"visits"`
	Users               int64   `json:"users"`
	BounceRate          float64 `json:"bounce_rate"`
	AvgVisitDurationSec int     `json:"avg_visit_duration_seconds"`
}

// AgeMetricsResult represents parsed age metrics result
type AgeMetricsResult struct {
	AgeGroup              string  `json:"age_group"`
	Visits                int64   `json:"visits"`
	Users                 int64   `json:"users"`
	BounceRate            float64 `json:"bounce_rate"`
	AvgSessionDurationSec int     `json:"avg_session_duration_seconds"`
}

// ConversionsResult represents parsed conversions result
type ConversionsResult struct {
	GoalID      int64 `json:"goal_id"`
	Visits      int64 `json:"visits"`
	Conversions int64 `json:"conversions"`
}

// Goal represents a goal from Management API
type Goal struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	IsRetargeting int    `json:"is_retargeting"` // API returns 0 or 1, not bool
}

// GoalsResponse represents response from Management API for goals
type GoalsResponse struct {
	Goals []Goal `json:"goals"`
}

// GetGoals retrieves list of goals for a counter
// Documentation: https://yandex.ru/dev/metrika/doc/api2/management_v1/goals.html
func (c *YandexMetricaClient) GetGoals(ctx context.Context, counterID int64) ([]Goal, error) {
	url := fmt.Sprintf("%s/counter/%d/goals", yandexMetricaManagementURL, counterID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Try to parse as array first (Management API may return array directly)
	var goalsArray []Goal
	if err := json.Unmarshal(body, &goalsArray); err == nil {
		return goalsArray, nil
	}

	// If array parsing failed, try as object with "goals" field
	var response GoalsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response (tried both array and object): %w, body: %s", err, string(body))
	}

	return response.Goals, nil
}

// GetMetrics retrieves metrics for a counter
// Documentation: https://yandex.ru/dev/metrika/doc/api2/api_v1/data.html
func (c *YandexMetricaClient) GetMetrics(ctx context.Context, counterID int64, dateFrom, dateTo string) (*MetricsResult, error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(counterID, 10))
	params.Set("date1", dateFrom)
	params.Set("date2", dateTo)
	params.Set("metrics", "ym:s:visits,ym:s:users,ym:s:bounceRate,ym:s:avgVisitDurationSeconds")
	params.Set("accuracy", "full")

	var response MetricsResponse
	if err := c.makeRequest(ctx, params, &response); err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	// Parse response
	if len(response.Data) == 0 || len(response.Data[0].Metrics) < 4 {
		return &MetricsResult{}, nil // Return empty result if no data
	}

	metrics := response.Data[0].Metrics
	return &MetricsResult{
		Visits:              int64(metrics[0]),
		Users:               int64(metrics[1]),
		BounceRate:          metrics[2],
		AvgVisitDurationSec: int(metrics[3]),
	}, nil
}

// GetMetricsByAge retrieves metrics broken down by age
// Documentation: https://yandex.ru/dev/metrika/doc/api2/api_v1/data.html
func (c *YandexMetricaClient) GetMetricsByAge(ctx context.Context, counterID int64, dateFrom, dateTo string) ([]AgeMetricsResult, error) {
	params := url.Values{}
	params.Set("ids", strconv.FormatInt(counterID, 10))
	params.Set("date1", dateFrom)
	params.Set("date2", dateTo)
	params.Set("metrics", "ym:s:visits,ym:s:users,ym:s:bounceRate,ym:s:avgVisitDurationSeconds")
	params.Set("dimensions", "ym:s:ageIntervalName")
	params.Set("accuracy", "full")

	var response AgeMetricsResponse
	if err := c.makeRequest(ctx, params, &response); err != nil {
		return nil, fmt.Errorf("failed to get age metrics: %w", err)
	}

	// Parse response
	results := make([]AgeMetricsResult, 0, len(response.Data))
	for _, row := range response.Data {
		if len(row.Dimensions) == 0 || len(row.Metrics) < 4 {
			continue
		}

		ageGroup := row.Dimensions[0].Name
		results = append(results, AgeMetricsResult{
			AgeGroup:              ageGroup,
			Visits:                int64(row.Metrics[0]),
			Users:                 int64(row.Metrics[1]),
			BounceRate:            row.Metrics[2],
			AvgSessionDurationSec: int(row.Metrics[3]),
		})
	}

	return results, nil
}

// GetConversions retrieves conversions for specified goals
// Documentation: https://yandex.ru/dev/metrika/doc/api2/api_v1/data.html
func (c *YandexMetricaClient) GetConversions(ctx context.Context, counterID int64, goalIDs []int64, dateFrom, dateTo string) ([]ConversionsResult, error) {
	if len(goalIDs) == 0 {
		return []ConversionsResult{}, nil
	}

	// Build metrics list for goals
	// Format: ym:s:goal<goal_id>visits (visits with goal) and ym:s:goal<goal_id>reaches (conversions)
	metricsList := make([]string, 0, len(goalIDs)*2)
	for _, goalID := range goalIDs {
		metricsList = append(metricsList,
			fmt.Sprintf("ym:s:goal%dvisits", goalID),
			fmt.Sprintf("ym:s:goal%dreaches", goalID),
		)
	}

	params := url.Values{}
	params.Set("ids", strconv.FormatInt(counterID, 10))
	params.Set("date1", dateFrom)
	params.Set("date2", dateTo)
	params.Set("metrics", strings.Join(metricsList, ","))
	params.Set("accuracy", "full")

	var response ConversionsResponse
	if err := c.makeRequest(ctx, params, &response); err != nil {
		return nil, fmt.Errorf("failed to get conversions: %w", err)
	}

	// Parse response
	results := make([]ConversionsResult, 0, len(goalIDs))
	if len(response.Data) == 0 || len(response.Data[0].Metrics) == 0 {
		// Return empty results for each goal
		for _, goalID := range goalIDs {
			results = append(results, ConversionsResult{
				GoalID:      goalID,
				Visits:      0,
				Conversions: 0,
			})
		}
		return results, nil
	}

	// Parse metrics: each goal has 2 metrics (visits with goal, reaches/conversions)
	metrics := response.Data[0].Metrics
	for i, goalID := range goalIDs {
		visitsIdx := i * 2        // ym:s:goal<id>visits
		conversionsIdx := i*2 + 1 // ym:s:goal<id>reaches

		visits := int64(0)
		conversions := int64(0)

		if visitsIdx < len(metrics) {
			visits = int64(metrics[visitsIdx])
		}
		if conversionsIdx < len(metrics) {
			conversions = int64(metrics[conversionsIdx])
		}

		results = append(results, ConversionsResult{
			GoalID:      goalID,
			Visits:      visits,
			Conversions: conversions,
		})
	}

	return results, nil
}

// makeRequest performs HTTP request to Yandex Metrica API
func (c *YandexMetricaClient) makeRequest(ctx context.Context, params url.Values, response interface{}) error {
	reqURL := c.baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers according to documentation
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

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
