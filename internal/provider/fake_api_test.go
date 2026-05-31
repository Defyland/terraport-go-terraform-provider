package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
)

type fakeBankPortAPI struct {
	t         *testing.T
	server    *httptest.Server
	token     string
	mu        sync.Mutex
	nextID    int
	counts    map[string]int
	fail429   map[string]int
	delays    map[string]time.Duration
	products  map[string]bankport.APIProduct
	apps      map[string]bankport.PartnerApp
	webhooks  map[string]bankport.WebhookEndpoint
	policies  map[string]bankport.RateLimitPolicy
	sandboxes map[string]bankport.SandboxEnvironment
}

func newFakeBankPortAPI(t *testing.T) *fakeBankPortAPI {
	t.Helper()
	api := &fakeBankPortAPI{
		t:         t,
		token:     "test-token",
		nextID:    1,
		counts:    map[string]int{},
		fail429:   map[string]int{},
		delays:    map[string]time.Duration{},
		products:  map[string]bankport.APIProduct{},
		apps:      map[string]bankport.PartnerApp{},
		webhooks:  map[string]bankport.WebhookEndpoint{},
		policies:  map[string]bankport.RateLimitPolicy{},
		sandboxes: map[string]bankport.SandboxEnvironment{},
	}
	api.products["bankport"] = bankport.APIProduct{
		Code:         "bankport",
		Name:         "BankPort Partner API",
		Category:     "payments",
		Beta:         false,
		Regions:      []string{"us-east-1", "sa-east-1"},
		Capabilities: []string{"partner-apps", "webhooks", "rate-limits", "sandbox-environments"},
		DocsURL:      "https://docs.bankport.example.test/products/bankport",
	}
	api.products["pixguard"] = bankport.APIProduct{
		Code:         "pixguard",
		Name:         "PixGuard Risk API",
		Category:     "risk",
		Beta:         true,
		Regions:      []string{"sa-east-1"},
		Capabilities: []string{"risk-scores", "fraud-signals"},
		DocsURL:      "https://docs.bankport.example.test/products/pixguard",
	}
	api.server = httptest.NewServer(http.HandlerFunc(api.handle))
	return api
}

func (f *fakeBankPortAPI) URL() string {
	return f.server.URL
}

func (f *fakeBankPortAPI) Close() {
	f.server.Close()
}

func (f *fakeBankPortAPI) Token() string {
	return f.token
}

func (f *fakeBankPortAPI) Fail429(methodPath string, count int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fail429[methodPath] = count
}

func (f *fakeBankPortAPI) Delay(methodPath string, delay time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.delays[methodPath] = delay
}

func (f *fakeBankPortAPI) RequestCount(methodPath string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.counts[methodPath]
}

func (f *fakeBankPortAPI) TotalRequests() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	total := 0
	for _, count := range f.counts {
		total += count
	}
	return total
}

func (f *fakeBankPortAPI) PartnerAppCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.apps)
}

func (f *fakeBankPortAPI) MutateFirstPartnerAppName(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for id, app := range f.apps {
		app.Name = name
		f.apps[id] = app
		return
	}
	f.t.Fatal("no partner app to mutate")
}

func (f *fakeBankPortAPI) handle(w http.ResponseWriter, r *http.Request) {
	key := r.Method + " " + r.URL.Path
	f.mu.Lock()
	f.counts[key]++
	delay := f.delays[key]
	if remaining := f.fail429[key]; remaining > 0 {
		f.fail429[key] = remaining - 1
		f.mu.Unlock()
		writeFakeError(w, http.StatusTooManyRequests, "rate_limited", "fake API rate limit")
		return
	}
	f.mu.Unlock()

	if delay > 0 {
		time.Sleep(delay)
	}

	if r.Header.Get("Authorization") != "Bearer "+f.token {
		writeFakeError(w, http.StatusUnauthorized, "unauthorized", "invalid bearer token")
		return
	}

	switch {
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/v1/products/"):
		f.handleProduct(w, r)
	case r.URL.Path == "/v1/partner-apps" && r.Method == http.MethodPost:
		f.createPartnerApp(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/partner-apps/"):
		f.handlePartnerApp(w, r)
	case r.URL.Path == "/v1/webhook-endpoints" && r.Method == http.MethodPost:
		f.createWebhookEndpoint(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/webhook-endpoints/"):
		f.handleWebhookEndpoint(w, r)
	case r.URL.Path == "/v1/rate-limit-policies" && r.Method == http.MethodPost:
		f.createRateLimitPolicy(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/rate-limit-policies/"):
		f.handleRateLimitPolicy(w, r)
	case r.URL.Path == "/v1/sandbox-environments" && r.Method == http.MethodPost:
		f.createSandboxEnvironment(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/sandbox-environments/"):
		f.handleSandboxEnvironment(w, r)
	default:
		writeFakeError(w, http.StatusNotFound, "not_found", "unknown fake API route")
	}
}

func (f *fakeBankPortAPI) handleProduct(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/v1/products/")
	f.mu.Lock()
	product, ok := f.products[code]
	f.mu.Unlock()
	if !ok {
		writeFakeError(w, http.StatusNotFound, "not_found", "product not found")
		return
	}
	writeFakeJSON(w, http.StatusOK, product)
}

func (f *fakeBankPortAPI) createPartnerApp(w http.ResponseWriter, r *http.Request) {
	var app bankport.PartnerApp
	decodeFakeJSON(w, r, &app)
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.newIDLocked("app")
	app.ID = id
	app.ClientID = "client_" + id
	app.ClientSecret = "client_secret_" + id
	if app.Status == "" {
		app.Status = "active"
	}
	f.apps[id] = app
	writeFakeJSON(w, http.StatusCreated, app)
}

func (f *fakeBankPortAPI) handlePartnerApp(w http.ResponseWriter, r *http.Request) {
	id, action := splitResourcePath(r.URL.Path, "/v1/partner-apps/")
	f.mu.Lock()
	app, ok := f.apps[id]
	f.mu.Unlock()
	if !ok {
		writeFakeError(w, http.StatusNotFound, "not_found", "partner app not found")
		return
	}
	if action == "rotate-secret" && r.Method == http.MethodPost {
		f.mu.Lock()
		app.ClientSecret = fmt.Sprintf("client_secret_%s_rotated_%d", id, f.nextID)
		f.nextID++
		f.apps[id] = app
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, bankport.SecretRotation{ClientSecret: app.ClientSecret})
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeFakeJSON(w, http.StatusOK, app)
	case http.MethodPut:
		var next bankport.PartnerApp
		decodeFakeJSON(w, r, &next)
		next.ID = app.ID
		next.ClientID = app.ClientID
		next.ClientSecret = app.ClientSecret
		f.mu.Lock()
		f.apps[id] = next
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, next)
	case http.MethodDelete:
		f.mu.Lock()
		delete(f.apps, id)
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	default:
		writeFakeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "unsupported method")
	}
}

func (f *fakeBankPortAPI) createWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpoint bankport.WebhookEndpoint
	decodeFakeJSON(w, r, &endpoint)
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.newIDLocked("wh")
	endpoint.ID = id
	endpoint.SigningSecret = "signing_secret_" + id
	f.webhooks[id] = endpoint
	writeFakeJSON(w, http.StatusCreated, endpoint)
}

func (f *fakeBankPortAPI) handleWebhookEndpoint(w http.ResponseWriter, r *http.Request) {
	id, action := splitResourcePath(r.URL.Path, "/v1/webhook-endpoints/")
	f.mu.Lock()
	endpoint, ok := f.webhooks[id]
	f.mu.Unlock()
	if !ok {
		writeFakeError(w, http.StatusNotFound, "not_found", "webhook endpoint not found")
		return
	}
	if action == "rotate-secret" && r.Method == http.MethodPost {
		f.mu.Lock()
		endpoint.SigningSecret = fmt.Sprintf("signing_secret_%s_rotated_%d", id, f.nextID)
		f.nextID++
		f.webhooks[id] = endpoint
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, bankport.SecretRotation{SigningSecret: endpoint.SigningSecret})
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeFakeJSON(w, http.StatusOK, endpoint)
	case http.MethodPut:
		var next bankport.WebhookEndpoint
		decodeFakeJSON(w, r, &next)
		next.ID = endpoint.ID
		next.SigningSecret = endpoint.SigningSecret
		f.mu.Lock()
		f.webhooks[id] = next
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, next)
	case http.MethodDelete:
		f.mu.Lock()
		delete(f.webhooks, id)
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	default:
		writeFakeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "unsupported method")
	}
}

func (f *fakeBankPortAPI) createRateLimitPolicy(w http.ResponseWriter, r *http.Request) {
	var policy bankport.RateLimitPolicy
	decodeFakeJSON(w, r, &policy)
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.newIDLocked("rlp")
	policy.ID = id
	if policy.Mode == "" {
		policy.Mode = "enforce"
	}
	f.policies[id] = policy
	writeFakeJSON(w, http.StatusCreated, policy)
}

func (f *fakeBankPortAPI) handleRateLimitPolicy(w http.ResponseWriter, r *http.Request) {
	id, _ := splitResourcePath(r.URL.Path, "/v1/rate-limit-policies/")
	f.mu.Lock()
	policy, ok := f.policies[id]
	f.mu.Unlock()
	if !ok {
		writeFakeError(w, http.StatusNotFound, "not_found", "rate-limit policy not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeFakeJSON(w, http.StatusOK, policy)
	case http.MethodPut:
		var next bankport.RateLimitPolicy
		decodeFakeJSON(w, r, &next)
		next.ID = policy.ID
		f.mu.Lock()
		f.policies[id] = next
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, next)
	case http.MethodDelete:
		f.mu.Lock()
		delete(f.policies, id)
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	default:
		writeFakeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "unsupported method")
	}
}

func (f *fakeBankPortAPI) createSandboxEnvironment(w http.ResponseWriter, r *http.Request) {
	var env bankport.SandboxEnvironment
	decodeFakeJSON(w, r, &env)
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.newIDLocked("sbx")
	env.ID = id
	env.APIKeyToken = "api_key_token_" + id
	if env.Region == "" {
		env.Region = "us-east-1"
	}
	if env.Status == "" {
		env.Status = "ready"
	}
	f.sandboxes[id] = env
	writeFakeJSON(w, http.StatusCreated, env)
}

func (f *fakeBankPortAPI) handleSandboxEnvironment(w http.ResponseWriter, r *http.Request) {
	id, _ := splitResourcePath(r.URL.Path, "/v1/sandbox-environments/")
	f.mu.Lock()
	env, ok := f.sandboxes[id]
	f.mu.Unlock()
	if !ok {
		writeFakeError(w, http.StatusNotFound, "not_found", "sandbox environment not found")
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeFakeJSON(w, http.StatusOK, env)
	case http.MethodPut:
		var next bankport.SandboxEnvironment
		decodeFakeJSON(w, r, &next)
		next.ID = env.ID
		next.APIKeyToken = env.APIKeyToken
		f.mu.Lock()
		f.sandboxes[id] = next
		f.mu.Unlock()
		writeFakeJSON(w, http.StatusOK, next)
	case http.MethodDelete:
		f.mu.Lock()
		delete(f.sandboxes, id)
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	default:
		writeFakeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "unsupported method")
	}
}

func (f *fakeBankPortAPI) newIDLocked(prefix string) string {
	id := fmt.Sprintf("%s_%04d", prefix, f.nextID)
	f.nextID++
	return id
}

func splitResourcePath(path, prefix string) (id string, action string) {
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.SplitN(rest, "/", 2)
	id = parts[0]
	if len(parts) == 2 {
		action = parts[1]
	}
	return id, action
}

func decodeFakeJSON(w http.ResponseWriter, r *http.Request, target any) {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeFakeError(w, http.StatusBadRequest, "invalid_json", "request body is not valid JSON")
	}
}

func writeFakeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeFakeError(w http.ResponseWriter, status int, code, message string) {
	writeFakeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
