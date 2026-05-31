package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultEndpoint = "https://api.terraport.local"

type terraportProvider struct {
	version string
}

type providerModel struct {
	Endpoint         types.String `tfsdk:"endpoint"`
	Token            types.String `tfsdk:"token"`
	TimeoutMS        types.Int64  `tfsdk:"timeout_ms"`
	RetryMaxAttempts types.Int64  `tfsdk:"retry_max_attempts"`
	RetryMinDelayMS  types.Int64  `tfsdk:"retry_min_delay_ms"`
}

var _ provider.Provider = (*terraportProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &terraportProvider{version: version}
	}
}

func (p *terraportProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "terraport"
	resp.Version = p.version
}

func (p *terraportProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraport configures fictitious BankPort platform resources through the Terraform Plugin Framework.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "BankPort-compatible API base URL. Falls back to `TERRAPORT_ENDPOINT`, then `BANKPORT_ENDPOINT`, then a documented placeholder endpoint.",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Bearer token for the BankPort API. Falls back to `TERRAPORT_TOKEN`, then `BANKPORT_TOKEN`.",
			},
			"timeout_ms": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "HTTP request timeout in milliseconds. Defaults to 10000.",
			},
			"retry_max_attempts": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Maximum attempts for retryable `429` and `5xx` responses. Defaults to 3.",
			},
			"retry_min_delay_ms": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Minimum exponential backoff delay in milliseconds. Defaults to 100.",
			},
		},
	}
}

func (p *terraportProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := valueOrEnv(config.Endpoint, "TERRAPORT_ENDPOINT", "BANKPORT_ENDPOINT")
	if endpoint == "" {
		endpoint = defaultEndpoint
	}
	token := valueOrEnv(config.Token, "TERRAPORT_TOKEN", "BANKPORT_TOKEN")
	timeoutMS := int64ValueOrDefault(config.TimeoutMS, 10000)
	retryMaxAttempts := int64ValueOrDefault(config.RetryMaxAttempts, 3)
	retryMinDelayMS := int64ValueOrDefault(config.RetryMinDelayMS, 100)

	client, err := bankport.NewClient(bankport.Config{
		Endpoint:    endpoint,
		Token:       token,
		Timeout:     time.Duration(timeoutMS) * time.Millisecond,
		MaxAttempts: int(retryMaxAttempts),
		MinBackoff:  time.Duration(retryMinDelayMS) * time.Millisecond,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Terraport provider configuration",
			fmt.Sprintf("Unable to configure BankPort API client: %s", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *terraportProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPartnerAppResource,
		NewWebhookEndpointResource,
		NewRateLimitPolicyResource,
		NewSandboxEnvironmentResource,
	}
}

func (p *terraportProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAPIProductDataSource,
	}
}

func valueOrEnv(value types.String, envNames ...string) string {
	if !value.IsNull() && !value.IsUnknown() {
		return value.ValueString()
	}
	for _, name := range envNames {
		if envValue := os.Getenv(name); envValue != "" {
			return envValue
		}
	}
	return ""
}

func int64ValueOrDefault(value types.Int64, fallback int64) int64 {
	if value.IsNull() || value.IsUnknown() || value.ValueInt64() <= 0 {
		return fallback
	}
	return value.ValueInt64()
}
