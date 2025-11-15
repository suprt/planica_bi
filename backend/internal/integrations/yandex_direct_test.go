package integrations

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// loadTestEnv loads .env file for testing
func loadTestEnv() {
	_ = godotenv.Load("../.env")      // Try backend/.env
	_ = godotenv.Load("../../.env")   // Try root .env
	_ = godotenv.Load(".env")         // Try current dir
}

// TestYandexDirectClient_RealAPI tests the client with real Yandex Direct API
// To run this test, ensure .env file contains:
// - YANDEX_OAUTH_TOKEN (OAuth token from OAuth flow)
// - YANDEX_DIRECT_TEST_CLIENT_LOGIN (your client login in Direct)
func TestYandexDirectClient_RealAPI(t *testing.T) {
	loadTestEnv()

	// Get token from environment
	token := os.Getenv("YANDEX_OAUTH_TOKEN")
	if token == "" {
		t.Skip("YANDEX_OAUTH_TOKEN not set in .env, skipping test")
	}

	clientLogin := os.Getenv("YANDEX_DIRECT_TEST_CLIENT_LOGIN")
	if clientLogin == "" {
		t.Skip("YANDEX_DIRECT_TEST_CLIENT_LOGIN not set in .env, skipping test")
	}

	// Create client with real API (not sandbox)
	client := NewYandexDirectClient(token, clientLogin, false)
	ctx := context.Background()

	// Test GetCampaigns
	t.Run("GetCampaigns", func(t *testing.T) {
		t.Logf("Testing GetCampaigns with Client-Login: %s", clientLogin)
		
		campaigns, err := client.GetCampaigns(ctx)
		if err != nil {
			t.Fatalf("GetCampaigns failed: %v", err)
		}

		t.Logf("✅ Successfully got %d campaigns", len(campaigns))
		if len(campaigns) == 0 {
			t.Log("⚠️  No campaigns found. Make sure you have at least one campaign in your Direct account.")
		}

		for i, campaign := range campaigns {
			t.Logf("  Campaign %d: ID=%d, Name=%s", i+1, campaign.Id, campaign.Name)
		}
	})

	// Test GetCampaignReport
	t.Run("GetCampaignReport", func(t *testing.T) {
		// Use recent date range (last month)
		dateFrom := "2024-12-01"
		dateTo := "2024-12-31"

		t.Logf("Testing GetCampaignReport from %s to %s", dateFrom, dateTo)

		report, err := client.GetCampaignReport(ctx, dateFrom, dateTo)
		if err != nil {
			t.Logf("⚠️  GetCampaignReport failed (may be expected if no data for this period): %v", err)
			return
		}

		t.Logf("✅ Successfully got %d report rows", len(report))
		if len(report) == 0 {
			t.Log("⚠️  No report data found for this period. Campaigns may not have data yet.")
		}

		for i, row := range report {
			t.Logf("  Report row %d: CampaignId=%d, CampaignName=%s, Impressions=%d, Clicks=%d, Cost=%.2f, CTR=%.2f%%",
				i+1, row.CampaignId, row.CampaignName, row.Impressions, row.Clicks, row.Cost, row.CTR)
		}
	})
}

// TestYandexDirectClient_Sandbox tests the client with sandbox environment
// To run this test, set environment variables:
// - YANDEX_OAUTH_TOKEN (OAuth token from .env)
// - YANDEX_DIRECT_TEST_CLIENT_LOGIN (test client login)
func TestYandexDirectClient_Sandbox(t *testing.T) {
	loadTestEnv()

	// Skip test if no token provided
	token := os.Getenv("YANDEX_OAUTH_TOKEN")
	if token == "" {
		t.Skip("YANDEX_OAUTH_TOKEN not set, skipping test")
	}

	clientLogin := os.Getenv("YANDEX_DIRECT_TEST_CLIENT_LOGIN")
	if clientLogin == "" {
		t.Skip("YANDEX_DIRECT_TEST_CLIENT_LOGIN not set, skipping test")
	}

	// Create client with sandbox enabled
	client := NewYandexDirectClient(token, clientLogin, true)
	ctx := context.Background()

	// Test GetCampaigns
	t.Run("GetCampaigns", func(t *testing.T) {
		campaigns, err := client.GetCampaigns(ctx)
		if err != nil {
			t.Logf("GetCampaigns error (may be expected in sandbox): %v", err)
			// Don't fail test, sandbox may not have campaigns
			return
		}

		t.Logf("Got %d campaigns", len(campaigns))
		for _, campaign := range campaigns {
			t.Logf("Campaign: ID=%d, Name=%s", campaign.Id, campaign.Name)
		}
	})

	// Test GetCampaignReport
	t.Run("GetCampaignReport", func(t *testing.T) {
		// Use recent date range for testing
		dateFrom := "2024-01-01"
		dateTo := "2024-01-31"

		report, err := client.GetCampaignReport(ctx, dateFrom, dateTo)
		if err != nil {
			t.Logf("GetCampaignReport error (may be expected in sandbox): %v", err)
			// Don't fail test, sandbox may not have report data
			return
		}

		t.Logf("Got %d report rows", len(report))
		for _, row := range report {
			t.Logf("Report row: CampaignId=%d, CampaignName=%s, Impressions=%d, Clicks=%d, Cost=%.2f",
				row.CampaignId, row.CampaignName, row.Impressions, row.Clicks, row.Cost)
		}
	})
}

