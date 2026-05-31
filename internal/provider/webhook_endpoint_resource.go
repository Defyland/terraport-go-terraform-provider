package provider

import (
	"context"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type webhookEndpointResource struct {
	client *bankport.Client
}

type webhookEndpointModel struct {
	ID                   types.String `tfsdk:"id"`
	PartnerAppID         types.String `tfsdk:"partner_app_id"`
	URL                  types.String `tfsdk:"url"`
	EventTypes           types.Set    `tfsdk:"event_types"`
	Enabled              types.Bool   `tfsdk:"enabled"`
	SigningSecret        types.String `tfsdk:"signing_secret"`
	SigningSecretVersion types.Int64  `tfsdk:"signing_secret_version"`
}

var _ resource.Resource = (*webhookEndpointResource)(nil)
var _ resource.ResourceWithConfigure = (*webhookEndpointResource)(nil)
var _ resource.ResourceWithImportState = (*webhookEndpointResource)(nil)

func NewWebhookEndpointResource() resource.Resource {
	return &webhookEndpointResource{}
}

func (r *webhookEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bankport_webhook_endpoint"
}

func (r *webhookEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provisions a BankPort webhook endpoint for partner events.",
		Attributes: map[string]schema.Attribute{
			"id":                     schema.StringAttribute{Computed: true},
			"partner_app_id":         schema.StringAttribute{Required: true},
			"url":                    schema.StringAttribute{Required: true},
			"event_types":            schema.SetAttribute{Required: true, ElementType: types.StringType},
			"enabled":                schema.BoolAttribute{Optional: true, Computed: true},
			"signing_secret":         schema.StringAttribute{Computed: true, Sensitive: true},
			"signing_secret_version": schema.Int64Attribute{Optional: true, Computed: true},
		},
	}
}

func (r *webhookEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (r *webhookEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := bankport.WebhookEndpoint{
		PartnerAppID: plan.PartnerAppID.ValueString(),
		URL:          plan.URL.ValueString(),
		EventTypes:   stringSetToSlice(ctx, plan.EventTypes, &resp.Diagnostics),
		Enabled:      boolValueOrDefault(plan.Enabled, true),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateWebhookEndpoint(ctx, endpoint)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to create BankPort webhook endpoint", err)
		return
	}

	plan.applyWebhookEndpoint(ctx, created, &resp.Diagnostics)
	if plan.SigningSecretVersion.IsNull() || plan.SigningSecretVersion.IsUnknown() {
		plan.SigningSecretVersion = types.Int64Value(1)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.GetWebhookEndpoint(ctx, state.ID.ValueString())
	if bankport.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to read BankPort webhook endpoint", err)
		return
	}

	secretVersion := state.SigningSecretVersion
	state.applyWebhookEndpoint(ctx, endpoint, &resp.Diagnostics)
	state.SigningSecretVersion = secretVersion
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *webhookEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookEndpointModel
	var state webhookEndpointModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := int64ValueOrDefault(plan.SigningSecretVersion, int64ValueOrDefault(state.SigningSecretVersion, 1))
	endpoint := bankport.WebhookEndpoint{
		PartnerAppID: plan.PartnerAppID.ValueString(),
		URL:          plan.URL.ValueString(),
		EventTypes:   stringSetToSlice(ctx, plan.EventTypes, &resp.Diagnostics),
		Enabled:      boolValueOrDefault(plan.Enabled, true),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateWebhookEndpoint(ctx, state.ID.ValueString(), endpoint)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to update BankPort webhook endpoint", err)
		return
	}
	if version > int64ValueOrDefault(state.SigningSecretVersion, 1) {
		rotated, err := r.client.RotateWebhookSigningSecret(ctx, state.ID.ValueString())
		if err != nil {
			addAPIError(&resp.Diagnostics, "Unable to rotate BankPort webhook signing secret", err)
			return
		}
		updated.SigningSecret = rotated.SigningSecret
	}

	plan.applyWebhookEndpoint(ctx, updated, &resp.Diagnostics)
	plan.SigningSecretVersion = types.Int64Value(version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *webhookEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookEndpointModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteWebhookEndpoint(ctx, state.ID.ValueString()); err != nil && !bankport.IsNotFound(err) {
		addAPIError(&resp.Diagnostics, "Unable to delete BankPort webhook endpoint", err)
	}
}

func (r *webhookEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *webhookEndpointModel) applyWebhookEndpoint(ctx context.Context, endpoint bankport.WebhookEndpoint, diags *diag.Diagnostics) {
	m.ID = types.StringValue(endpoint.ID)
	m.PartnerAppID = types.StringValue(endpoint.PartnerAppID)
	m.URL = types.StringValue(endpoint.URL)
	m.EventTypes = stringSetValue(ctx, endpoint.EventTypes, diags)
	m.Enabled = types.BoolValue(endpoint.Enabled)
	m.SigningSecret = types.StringValue(endpoint.SigningSecret)
}
