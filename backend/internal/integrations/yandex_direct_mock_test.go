package integrations

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestYandexDirectClient_GetCampaigns_Mock tests GetCampaigns with mocked HTTP server
func TestYandexDirectClient_GetCampaigns_Mock(t *testing.T) {
	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Authorization") != "Bearer test_token" {
			t.Errorf("Expected Authorization header 'Bearer test_token', got '%s'", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Client-Login") != "test_client_login" {
			t.Errorf("Expected Client-Login header 'test_client_login', got '%s'", r.Header.Get("Client-Login"))
		}
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", r.Header.Get("Content-Type"))
		}

		// Verify request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if requestBody["method"] != "get" {
			t.Errorf("Expected method 'get', got '%v'", requestBody["method"])
		}

		// Return mock response
		response := CampaignsResponse{
			Result: struct {
				Campaigns []Campaign `json:"Campaigns"`
			}{
				Campaigns: []Campaign{
					{Id: 1, Name: "Test Campaign 1"},
					{Id: 2, Name: "Test Campaign 2"},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create client with mock server URL
	client := NewYandexDirectClientWithURL("test_token", "test_client_login", mockServer.URL)

	ctx := context.Background()

	// Test GetCampaigns
	campaigns, err := client.GetCampaigns(ctx)
	if err != nil {
		t.Fatalf("GetCampaigns failed: %v", err)
	}

	if len(campaigns) != 2 {
		t.Errorf("Expected 2 campaigns, got %d", len(campaigns))
	}

	if campaigns[0].Id != 1 || campaigns[0].Name != "Test Campaign 1" {
		t.Errorf("Expected campaign {Id: 1, Name: 'Test Campaign 1'}, got {Id: %d, Name: '%s'}", campaigns[0].Id, campaigns[0].Name)
	}
}

// TestYandexDirectClient_GetCampaignReport_Mock tests GetCampaignReport with mocked HTTP server
func TestYandexDirectClient_GetCampaignReport_Mock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/reports" {
			t.Errorf("Expected path '/reports', got '%s'", r.URL.Path)
		}

		// Return mock report response
		response := ReportResponse{
			Result: struct {
				Report struct {
					Rows []ReportRow `json:"Rows"`
				} `json:"Report"`
			}{
				Report: struct {
					Rows []ReportRow `json:"Rows"`
				}{
					Rows: []ReportRow{
						{
							CampaignId:   1,
							CampaignName: "Test Campaign",
							Impressions:  1000,
							Clicks:       50,
							Cost:         500.0,
							CTR:          5.0,
							AvgCpc:       10.0,
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	client := NewYandexDirectClientWithURL("test_token", "test_client_login", mockServer.URL)
	ctx := context.Background()

	report, err := client.GetCampaignReport(ctx, "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetCampaignReport failed: %v", err)
	}

	if len(report) != 1 {
		t.Errorf("Expected 1 report row, got %d", len(report))
	}

	if report[0].Impressions != 1000 || report[0].Clicks != 50 {
		t.Errorf("Expected {Impressions: 1000, Clicks: 50}, got {Impressions: %d, Clicks: %d}", report[0].Impressions, report[0].Clicks)
	}
}

// TestYandexDirectClient_ErrorHandling tests error handling
func TestYandexDirectClient_ErrorHandling(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResponse := CampaignsResponse{
			Error: &APIError{
				ErrorString: "Invalid token",
				ErrorCode:   53,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // API returns 200 even on errors
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer mockServer.Close()

	client := NewYandexDirectClientWithURL("invalid_token", "test_client_login", mockServer.URL)
	ctx := context.Background()

	_, err := client.GetCampaigns(ctx)
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}

	if err.Error() != "API error: Invalid token (code: 53)" {
		t.Errorf("Expected error message 'API error: Invalid token (code: 53)', got '%s'", err.Error())
	}
}

