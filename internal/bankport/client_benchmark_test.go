package bankport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkClientCreate100PartnerApps(b *testing.B) {
	var nextID atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/partner-apps" {
			writeClientTestError(w, http.StatusNotFound, "not_found", "unexpected benchmark route")
			return
		}
		var app PartnerApp
		if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
			writeClientTestError(w, http.StatusBadRequest, "invalid_json", "invalid JSON")
			return
		}
		id := "app_bench_" + strconv.FormatInt(nextID.Add(1), 10)
		app.ID = id
		app.ClientID = "client_" + id
		app.ClientSecret = "client_secret_" + id
		_ = json.NewEncoder(w).Encode(app)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Endpoint:    server.URL,
		Token:       "bench-token",
		Timeout:     5 * time.Second,
		MaxAttempts: 1,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			_, err := client.CreatePartnerApp(ctx, PartnerApp{
				Name:         "Load App " + strconv.Itoa(i),
				ProductCode:  "bankport",
				RedirectURIs: []string{"https://partner.example.test/callback"},
				Scopes:       []string{"accounts:read"},
				Status:       "active",
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkClientRetry429Twice(b *testing.B) {
	var attempts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := attempts.Add(1)
		if current%3 != 0 {
			writeClientTestError(w, http.StatusTooManyRequests, "rate_limited", "benchmark retry")
			return
		}
		_ = json.NewEncoder(w).Encode(APIProduct{Code: "bankport", Name: "BankPort Partner API"})
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Endpoint:    server.URL,
		Token:       "bench-token",
		Timeout:     5 * time.Second,
		MaxAttempts: 3,
		MinBackoff:  time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := client.GetAPIProduct(ctx, "bankport"); err != nil {
			b.Fatal(err)
		}
	}
}
