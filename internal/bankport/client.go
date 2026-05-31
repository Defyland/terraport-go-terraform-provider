package bankport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const defaultUserAgent = "terraport-terraform-provider/0.1"

type Config struct {
	Endpoint    string
	Token       string
	Timeout     time.Duration
	MaxAttempts int
	MinBackoff  time.Duration
	HTTPClient  *http.Client
	UserAgent   string
}

type Client struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	maxAttempts int
	minBackoff  time.Duration
	userAgent   string
	metrics     Metrics
}

type Metrics struct {
	Requests             atomic.Int64
	Retries              atomic.Int64
	RateLimitResponses   atomic.Int64
	ServerErrorResponses atomic.Int64
}

type MetricsSnapshot struct {
	Requests             int64
	Retries              int64
	RateLimitResponses   int64
	ServerErrorResponses int64
}

func NewClient(cfg Config) (*Client, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, errors.New("bankport endpoint is required")
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse bankport endpoint: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("bankport endpoint must include scheme and host")
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	maxAttempts := cfg.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	minBackoff := cfg.MinBackoff
	if minBackoff <= 0 {
		minBackoff = 100 * time.Millisecond
	}
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}
	userAgent := strings.TrimSpace(cfg.UserAgent)
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	return &Client{
		baseURL:     strings.TrimRight(parsed.String(), "/"),
		token:       cfg.Token,
		httpClient:  httpClient,
		maxAttempts: maxAttempts,
		minBackoff:  minBackoff,
		userAgent:   userAgent,
	}, nil
}

func (c *Client) SnapshotMetrics() MetricsSnapshot {
	return MetricsSnapshot{
		Requests:             c.metrics.Requests.Load(),
		Retries:              c.metrics.Retries.Load(),
		RateLimitResponses:   c.metrics.RateLimitResponses.Load(),
		ServerErrorResponses: c.metrics.ServerErrorResponses.Load(),
	}
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Retryable  bool
}

func (e *APIError) Error() string {
	if e.Code == "" {
		return fmt.Sprintf("bankport api request failed: status=%d", e.StatusCode)
	}
	message := strings.TrimSpace(Redact(e.Message))
	if message == "" {
		return fmt.Sprintf("bankport api request failed: status=%d code=%s", e.StatusCode, e.Code)
	}
	return fmt.Sprintf("bankport api request failed: status=%d code=%s message=%s", e.StatusCode, e.Code, message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

var sensitiveValuePattern = regexp.MustCompile(`(?i)(token|secret|api_key|client_secret|signing_secret)(\s*[:=]\s*)("[^"]+"|'[^']+'|[^\s,}]+)`)

func Redact(input string) string {
	return sensitiveValuePattern.ReplaceAllString(input, `${1}${2}[REDACTED]`)
}

type apiErrorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload any, out any) error {
	var body []byte
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
		body = encoded
	}

	attempts := c.maxAttempts
	if attempts < 1 {
		attempts = 1
	}

	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("build bankport request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.userAgent)
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		c.metrics.Requests.Add(1)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("bankport api request failed: %w", err)
		}

		err = c.decodeResponse(ctx, resp, out)
		if err == nil {
			return nil
		}

		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Retryable && attempt < attempts {
			c.metrics.Retries.Add(1)
			if apiErr.StatusCode == http.StatusTooManyRequests {
				c.metrics.RateLimitResponses.Add(1)
			}
			if apiErr.StatusCode >= 500 {
				c.metrics.ServerErrorResponses.Add(1)
			}
			if sleepErr := sleepWithBackoff(ctx, resp.Header.Get("Retry-After"), c.minBackoff, attempt); sleepErr != nil {
				return sleepErr
			}
			continue
		}

		return err
	}

	return errors.New("bankport api request failed after retry attempts")
}

func (c *Client) decodeResponse(ctx context.Context, resp *http.Response, out any) error {
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read bankport response: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if out == nil || len(body) == 0 {
			return nil
		}
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode bankport response: %w", err)
		}
		return nil
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Retryable:  resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500,
	}
	if len(body) > 0 {
		var envelope apiErrorEnvelope
		if err := json.Unmarshal(body, &envelope); err == nil {
			apiErr.Code = envelope.Error.Code
			apiErr.Message = envelope.Error.Message
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return apiErr
}

func sleepWithBackoff(ctx context.Context, retryAfter string, minBackoff time.Duration, attempt int) error {
	delay := retryAfterDelay(retryAfter)
	if delay <= 0 {
		delay = minBackoff * time.Duration(1<<max(attempt-1, 0))
	}
	if delay > 2*time.Second {
		delay = 2 * time.Second
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func retryAfterDelay(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(value); err == nil {
		return time.Until(when)
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *Client) GetAPIProduct(ctx context.Context, code string) (APIProduct, error) {
	var product APIProduct
	err := c.doJSON(ctx, http.MethodGet, "/v1/products/"+url.PathEscape(code), nil, &product)
	return product, err
}

func (c *Client) CreatePartnerApp(ctx context.Context, app PartnerApp) (PartnerApp, error) {
	var created PartnerApp
	err := c.doJSON(ctx, http.MethodPost, "/v1/partner-apps", app, &created)
	return created, err
}

func (c *Client) GetPartnerApp(ctx context.Context, id string) (PartnerApp, error) {
	var app PartnerApp
	err := c.doJSON(ctx, http.MethodGet, "/v1/partner-apps/"+url.PathEscape(id), nil, &app)
	return app, err
}

func (c *Client) UpdatePartnerApp(ctx context.Context, id string, app PartnerApp) (PartnerApp, error) {
	var updated PartnerApp
	err := c.doJSON(ctx, http.MethodPut, "/v1/partner-apps/"+url.PathEscape(id), app, &updated)
	return updated, err
}

func (c *Client) DeletePartnerApp(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, "/v1/partner-apps/"+url.PathEscape(id), nil, nil)
}

func (c *Client) RotatePartnerAppSecret(ctx context.Context, id string) (SecretRotation, error) {
	var rotated SecretRotation
	err := c.doJSON(ctx, http.MethodPost, "/v1/partner-apps/"+url.PathEscape(id)+"/rotate-secret", nil, &rotated)
	return rotated, err
}

func (c *Client) CreateWebhookEndpoint(ctx context.Context, endpoint WebhookEndpoint) (WebhookEndpoint, error) {
	var created WebhookEndpoint
	err := c.doJSON(ctx, http.MethodPost, "/v1/webhook-endpoints", endpoint, &created)
	return created, err
}

func (c *Client) GetWebhookEndpoint(ctx context.Context, id string) (WebhookEndpoint, error) {
	var endpoint WebhookEndpoint
	err := c.doJSON(ctx, http.MethodGet, "/v1/webhook-endpoints/"+url.PathEscape(id), nil, &endpoint)
	return endpoint, err
}

func (c *Client) UpdateWebhookEndpoint(ctx context.Context, id string, endpoint WebhookEndpoint) (WebhookEndpoint, error) {
	var updated WebhookEndpoint
	err := c.doJSON(ctx, http.MethodPut, "/v1/webhook-endpoints/"+url.PathEscape(id), endpoint, &updated)
	return updated, err
}

func (c *Client) DeleteWebhookEndpoint(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, "/v1/webhook-endpoints/"+url.PathEscape(id), nil, nil)
}

func (c *Client) RotateWebhookSigningSecret(ctx context.Context, id string) (SecretRotation, error) {
	var rotated SecretRotation
	err := c.doJSON(ctx, http.MethodPost, "/v1/webhook-endpoints/"+url.PathEscape(id)+"/rotate-secret", nil, &rotated)
	return rotated, err
}

func (c *Client) CreateRateLimitPolicy(ctx context.Context, policy RateLimitPolicy) (RateLimitPolicy, error) {
	var created RateLimitPolicy
	err := c.doJSON(ctx, http.MethodPost, "/v1/rate-limit-policies", policy, &created)
	return created, err
}

func (c *Client) GetRateLimitPolicy(ctx context.Context, id string) (RateLimitPolicy, error) {
	var policy RateLimitPolicy
	err := c.doJSON(ctx, http.MethodGet, "/v1/rate-limit-policies/"+url.PathEscape(id), nil, &policy)
	return policy, err
}

func (c *Client) UpdateRateLimitPolicy(ctx context.Context, id string, policy RateLimitPolicy) (RateLimitPolicy, error) {
	var updated RateLimitPolicy
	err := c.doJSON(ctx, http.MethodPut, "/v1/rate-limit-policies/"+url.PathEscape(id), policy, &updated)
	return updated, err
}

func (c *Client) DeleteRateLimitPolicy(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, "/v1/rate-limit-policies/"+url.PathEscape(id), nil, nil)
}

func (c *Client) CreateSandboxEnvironment(ctx context.Context, env SandboxEnvironment) (SandboxEnvironment, error) {
	var created SandboxEnvironment
	err := c.doJSON(ctx, http.MethodPost, "/v1/sandbox-environments", env, &created)
	return created, err
}

func (c *Client) GetSandboxEnvironment(ctx context.Context, id string) (SandboxEnvironment, error) {
	var env SandboxEnvironment
	err := c.doJSON(ctx, http.MethodGet, "/v1/sandbox-environments/"+url.PathEscape(id), nil, &env)
	return env, err
}

func (c *Client) UpdateSandboxEnvironment(ctx context.Context, id string, env SandboxEnvironment) (SandboxEnvironment, error) {
	var updated SandboxEnvironment
	err := c.doJSON(ctx, http.MethodPut, "/v1/sandbox-environments/"+url.PathEscape(id), env, &updated)
	return updated, err
}

func (c *Client) DeleteSandboxEnvironment(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, "/v1/sandbox-environments/"+url.PathEscape(id), nil, nil)
}
