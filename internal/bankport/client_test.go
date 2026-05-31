package bankport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientRetriesRateLimitThenSucceeds(t *testing.T) {
	var attempts atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected authorization header %q", got)
		}
		if attempts.Add(1) < 3 {
			writeClientTestError(w, http.StatusTooManyRequests, "rate_limited", "retry later")
			return
		}
		_ = json.NewEncoder(w).Encode(APIProduct{
			Code:         "bankport",
			Name:         "BankPort Partner API",
			Category:     "payments",
			Regions:      []string{"us-east-1"},
			Capabilities: []string{"partner-apps"},
			DocsURL:      "https://docs.example.test/bankport",
		})
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Endpoint:    server.URL,
		Token:       "test-token",
		Timeout:     time.Second,
		MaxAttempts: 3,
		MinBackoff:  time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	product, err := client.GetAPIProduct(context.Background(), "bankport")
	if err != nil {
		t.Fatal(err)
	}
	if product.Code != "bankport" {
		t.Fatalf("unexpected product code %q", product.Code)
	}
	metrics := client.SnapshotMetrics()
	if metrics.Requests != 3 || metrics.Retries != 2 || metrics.RateLimitResponses != 2 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
}

func TestClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_ = json.NewEncoder(w).Encode(APIProduct{Code: "bankport"})
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Endpoint: server.URL,
		Token:    "test-token",
		Timeout:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetAPIProduct(context.Background(), "bankport")
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Fatalf("expected timeout-like error, got %q", err.Error())
	}
}

func TestAPIErrorRedactsSensitiveValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeClientTestError(w, http.StatusUnauthorized, "unauthorized", "token=plain-token client_secret=plain-secret signing_secret=webhook-secret")
	}))
	defer server.Close()

	client, err := NewClient(Config{Endpoint: server.URL, Token: "plain-token"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetAPIProduct(context.Background(), "bankport")
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
	for _, leaked := range []string{"plain-token", "plain-secret", "webhook-secret"} {
		if strings.Contains(err.Error(), leaked) {
			t.Fatalf("error leaked sensitive value %q: %s", leaked, err)
		}
	}
}

func writeClientTestError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
