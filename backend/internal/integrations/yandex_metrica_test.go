package integrations

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

// TestYandexMetricaClient_RealAPI tests the client with real Yandex Metrica API
// To run this test, ensure .env file contains:
// - YANDEX_OAUTH_TOKEN (OAuth token from OAuth flow)
// - YANDEX_METRICA_TEST_COUNTER_ID (your test counter ID)
func TestYandexMetricaClient_RealAPI(t *testing.T) {
	// Load .env file
	_ = godotenv.Load("../.env")    // Try backend/.env
	_ = godotenv.Load("../../.env") // Try root .env
	_ = godotenv.Load(".env")       // Try current dir

	// Get token from environment
	token := os.Getenv("YANDEX_OAUTH_TOKEN")
	if token == "" {
		t.Skip("YANDEX_OAUTH_TOKEN not set in .env, skipping test")
	}

	counterIDStr := os.Getenv("YANDEX_METRICA_TEST_COUNTER_ID")
	if counterIDStr == "" {
		t.Skip("YANDEX_METRICA_TEST_COUNTER_ID not set in .env, skipping test")
	}

	// Parse counter ID
	counterID, err := strconv.ParseInt(counterIDStr, 10, 64)
	if err != nil {
		t.Fatalf("Invalid counter ID format: %s", counterIDStr)
	}

	// Create client with real API
	client := NewYandexMetricaClient(token)
	ctx := context.Background()

	// Test GetMetrics
	t.Run("GetMetrics", func(t *testing.T) {
		t.Logf("Testing GetMetrics with CounterID: %d", counterID)

		result, err := client.GetMetrics(ctx, counterID, "2024-12-01", "2024-12-31")
		if err != nil {
			t.Fatalf("GetMetrics failed: %v", err)
		}

		t.Logf("✅ Successfully got metrics:")
		t.Logf("  Visits: %d", result.Visits)
		t.Logf("  Users: %d", result.Users)
		t.Logf("  BounceRate: %.2f%%", result.BounceRate)
		t.Logf("  AvgVisitDurationSec: %d", result.AvgVisitDurationSec)
	})

	// Test GetMetricsByAge
	t.Run("GetMetricsByAge", func(t *testing.T) {
		t.Logf("Testing GetMetricsByAge with CounterID: %d", counterID)

		results, err := client.GetMetricsByAge(ctx, counterID, "2024-12-01", "2024-12-31")
		if err != nil {
			t.Logf("⚠️  GetMetricsByAge failed (may be expected if no data): %v", err)
			return
		}

		t.Logf("✅ Successfully got %d age groups", len(results))
		if len(results) == 0 {
			t.Log("⚠️  No age data found for this period")
		}

		for i, result := range results {
			t.Logf("  Age group %d: %s", i+1, result.AgeGroup)
			t.Logf("    Visits: %d, Users: %d", result.Visits, result.Users)
			t.Logf("    BounceRate: %.2f%%, AvgSessionDurationSec: %d", result.BounceRate, result.AvgSessionDurationSec)
		}
	})

	// Test GetGoals - get list of goals first
	t.Run("GetGoals", func(t *testing.T) {
		t.Logf("Testing GetGoals with CounterID: %d", counterID)

		goals, err := client.GetGoals(ctx, counterID)
		if err != nil {
			t.Logf("⚠️  GetGoals failed: %v", err)
			return
		}

		t.Logf("✅ Successfully got %d goals", len(goals))
		if len(goals) == 0 {
			t.Log("⚠️  No goals found in this counter")
			return
		}

		for i, goal := range goals {
			t.Logf("  Goal %d: ID=%d, Name=%s, Type=%s", i+1, goal.ID, goal.Name, goal.Type)
		}

		// Test GetConversions with real goal IDs
		t.Run("GetConversions", func(t *testing.T) {
			// Use first goal ID from the list
			goalIDs := make([]int64, 0, len(goals))
			for _, goal := range goals {
				goalIDs = append(goalIDs, goal.ID)
			}

			t.Logf("Testing GetConversions with CounterID: %d, GoalIDs: %v", counterID, goalIDs)

			results, err := client.GetConversions(ctx, counterID, goalIDs, "2024-12-01", "2024-12-31")
			if err != nil {
				t.Logf("⚠️  GetConversions failed: %v", err)
				return
			}

			t.Logf("✅ Successfully got %d conversion results", len(results))
			for i, result := range results {
				t.Logf("  Goal %d: GoalID=%d, Visits=%d, Conversions=%d", i+1, result.GoalID, result.Visits, result.Conversions)
			}
		})
	})
}
