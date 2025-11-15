package integrations

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestYandexMetricaClient_GetMetrics_Mock tests GetMetrics with mocked HTTP server
func TestYandexMetricaClient_GetMetrics_Mock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "OAuth test_token" {
			t.Errorf("Expected Authorization header 'OAuth test_token', got '%s'", r.Header.Get("Authorization"))
		}

		// Verify query parameters
		if !r.URL.Query().Has("ids") {
			t.Error("Expected 'ids' query parameter")
		}
		if r.URL.Query().Get("ids") != "12345" {
			t.Errorf("Expected ids=12345, got ids=%s", r.URL.Query().Get("ids"))
		}

		// Return mock response
		response := MetricsResponse{
			Data: []MetricsData{
				{
					Dimensions: []Dimension{},
					Metrics:    []float64{1000, 800, 45.5, 120},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	client := NewYandexMetricaClientWithURL("test_token", mockServer.URL)
	ctx := context.Background()

	result, err := client.GetMetrics(ctx, 12345, "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetMetrics failed: %v", err)
	}

	if result.Visits != 1000 {
		t.Errorf("Expected Visits=1000, got %d", result.Visits)
	}
	if result.Users != 800 {
		t.Errorf("Expected Users=800, got %d", result.Users)
	}
	if result.BounceRate != 45.5 {
		t.Errorf("Expected BounceRate=45.5, got %.2f", result.BounceRate)
	}
	if result.AvgVisitDurationSec != 120 {
		t.Errorf("Expected AvgVisitDurationSec=120, got %d", result.AvgVisitDurationSec)
	}
}

// TestYandexMetricaClient_GetMetricsByAge_Mock tests GetMetricsByAge with mocked HTTP server
func TestYandexMetricaClient_GetMetricsByAge_Mock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("dimensions") != "ym:s:ageIntervalName" {
			t.Errorf("Expected dimensions=ym:s:ageIntervalName, got dimensions=%s", r.URL.Query().Get("dimensions"))
		}

		// Return mock response
		response := AgeMetricsResponse{
			Data: []AgeMetricsData{
				{
					Dimensions: []Dimension{{Name: "18-24"}},
					Metrics:    []float64{300, 250},
				},
				{
					Dimensions: []Dimension{{Name: "25-34"}},
					Metrics:    []float64{400, 320},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	client := NewYandexMetricaClientWithURL("test_token", mockServer.URL)
	ctx := context.Background()

	results, err := client.GetMetricsByAge(ctx, 12345, "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetMetricsByAge failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 age groups, got %d", len(results))
	}

	if results[0].AgeGroup != "18-24" || results[0].Visits != 300 {
		t.Errorf("Expected {AgeGroup: '18-24', Visits: 300}, got {AgeGroup: '%s', Visits: %d}",
			results[0].AgeGroup, results[0].Visits)
	}
}

// TestYandexMetricaClient_GetConversions_Mock tests GetConversions with mocked HTTP server
func TestYandexMetricaClient_GetConversions_Mock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify metrics parameter contains goal metrics
		metrics := r.URL.Query().Get("metrics")
		if !contains(metrics, "ym:s:goal123Visits") {
			t.Errorf("Expected metrics to contain 'ym:s:goal123Visits', got '%s'", metrics)
		}

		// Return mock response
		response := ConversionsResponse{
			Data: []ConversionsData{
				{
					Dimensions: []Dimension{},
					Metrics:    []float64{50, 10}, // goal123: visits=50, conversions=10
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	client := NewYandexMetricaClientWithURL("test_token", mockServer.URL)
	ctx := context.Background()

	results, err := client.GetConversions(ctx, 12345, []int64{123}, "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetConversions failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 conversion result, got %d", len(results))
	}

	if results[0].GoalID != 123 || results[0].Visits != 50 || results[0].Conversions != 10 {
		t.Errorf("Expected {GoalID: 123, Visits: 50, Conversions: 10}, got {GoalID: %d, Visits: %d, Conversions: %d}",
			results[0].GoalID, results[0].Visits, results[0].Conversions)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || strings.Contains(s, substr))
}
