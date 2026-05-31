package bankport

type PartnerApp struct {
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	ProductCode  string   `json:"product_code"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       []string `json:"scopes"`
	Status       string   `json:"status"`
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
}

type WebhookEndpoint struct {
	ID            string   `json:"id,omitempty"`
	PartnerAppID  string   `json:"partner_app_id"`
	URL           string   `json:"url"`
	EventTypes    []string `json:"event_types"`
	Enabled       bool     `json:"enabled"`
	SigningSecret string   `json:"signing_secret,omitempty"`
}

type RateLimitPolicy struct {
	ID                string `json:"id,omitempty"`
	ProductCode       string `json:"product_code"`
	SubjectType       string `json:"subject_type"`
	SubjectID         string `json:"subject_id"`
	RequestsPerMinute int64  `json:"requests_per_minute"`
	BurstLimit        int64  `json:"burst_limit"`
	Mode              string `json:"mode"`
}

type SandboxEnvironment struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name"`
	Products    []string `json:"products"`
	Region      string   `json:"region"`
	Status      string   `json:"status"`
	APIKeyToken string   `json:"api_key_token,omitempty"`
}

type APIProduct struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	Beta         bool     `json:"beta"`
	Regions      []string `json:"regions"`
	Capabilities []string `json:"capabilities"`
	DocsURL      string   `json:"docs_url"`
}

type SecretRotation struct {
	ClientSecret  string `json:"client_secret,omitempty"`
	SigningSecret string `json:"signing_secret,omitempty"`
}
