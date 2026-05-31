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

type partnerAppResource struct {
	client *bankport.Client
}

type partnerAppModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ProductCode         types.String `tfsdk:"product_code"`
	RedirectURIs        types.Set    `tfsdk:"redirect_uris"`
	Scopes              types.Set    `tfsdk:"scopes"`
	Status              types.String `tfsdk:"status"`
	ClientID            types.String `tfsdk:"client_id"`
	ClientSecret        types.String `tfsdk:"client_secret"`
	ClientSecretVersion types.Int64  `tfsdk:"client_secret_version"`
}

var _ resource.Resource = (*partnerAppResource)(nil)
var _ resource.ResourceWithConfigure = (*partnerAppResource)(nil)
var _ resource.ResourceWithImportState = (*partnerAppResource)(nil)

func NewPartnerAppResource() resource.Resource {
	return &partnerAppResource{}
}

func (r *partnerAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bankport_partner_app"
}

func (r *partnerAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provisions a BankPort partner OAuth application.",
		Attributes: map[string]schema.Attribute{
			"id":                    schema.StringAttribute{Computed: true},
			"name":                  schema.StringAttribute{Required: true},
			"product_code":          schema.StringAttribute{Required: true},
			"redirect_uris":         schema.SetAttribute{Required: true, ElementType: types.StringType},
			"scopes":                schema.SetAttribute{Required: true, ElementType: types.StringType},
			"status":                schema.StringAttribute{Optional: true, Computed: true},
			"client_id":             schema.StringAttribute{Computed: true},
			"client_secret":         schema.StringAttribute{Computed: true, Sensitive: true},
			"client_secret_version": schema.Int64Attribute{Optional: true, Computed: true},
		},
	}
}

func (r *partnerAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (r *partnerAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app := bankport.PartnerApp{
		Name:         plan.Name.ValueString(),
		ProductCode:  plan.ProductCode.ValueString(),
		RedirectURIs: stringSetToSlice(ctx, plan.RedirectURIs, &resp.Diagnostics),
		Scopes:       stringSetToSlice(ctx, plan.Scopes, &resp.Diagnostics),
		Status:       stringValueOrDefault(plan.Status, "active"),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreatePartnerApp(ctx, app)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to create BankPort partner app", err)
		return
	}

	plan.applyPartnerApp(ctx, created, &resp.Diagnostics)
	if plan.ClientSecretVersion.IsNull() || plan.ClientSecretVersion.IsUnknown() {
		plan.ClientSecretVersion = types.Int64Value(1)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *partnerAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.GetPartnerApp(ctx, state.ID.ValueString())
	if bankport.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to read BankPort partner app", err)
		return
	}

	secretVersion := state.ClientSecretVersion
	state.applyPartnerApp(ctx, app, &resp.Diagnostics)
	state.ClientSecretVersion = secretVersion
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *partnerAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan partnerAppModel
	var state partnerAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version := int64ValueOrDefault(plan.ClientSecretVersion, int64ValueOrDefault(state.ClientSecretVersion, 1))
	app := bankport.PartnerApp{
		Name:         plan.Name.ValueString(),
		ProductCode:  plan.ProductCode.ValueString(),
		RedirectURIs: stringSetToSlice(ctx, plan.RedirectURIs, &resp.Diagnostics),
		Scopes:       stringSetToSlice(ctx, plan.Scopes, &resp.Diagnostics),
		Status:       stringValueOrDefault(plan.Status, "active"),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdatePartnerApp(ctx, state.ID.ValueString(), app)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to update BankPort partner app", err)
		return
	}
	if version > int64ValueOrDefault(state.ClientSecretVersion, 1) {
		rotated, err := r.client.RotatePartnerAppSecret(ctx, state.ID.ValueString())
		if err != nil {
			addAPIError(&resp.Diagnostics, "Unable to rotate BankPort partner app client secret", err)
			return
		}
		updated.ClientSecret = rotated.ClientSecret
	}

	plan.applyPartnerApp(ctx, updated, &resp.Diagnostics)
	plan.ClientSecretVersion = types.Int64Value(version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *partnerAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state partnerAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeletePartnerApp(ctx, state.ID.ValueString()); err != nil && !bankport.IsNotFound(err) {
		addAPIError(&resp.Diagnostics, "Unable to delete BankPort partner app", err)
	}
}

func (r *partnerAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *partnerAppModel) applyPartnerApp(ctx context.Context, app bankport.PartnerApp, diags *diag.Diagnostics) {
	m.ID = types.StringValue(app.ID)
	m.Name = types.StringValue(app.Name)
	m.ProductCode = types.StringValue(app.ProductCode)
	m.RedirectURIs = stringSetValue(ctx, app.RedirectURIs, diags)
	m.Scopes = stringSetValue(ctx, app.Scopes, diags)
	m.Status = types.StringValue(app.Status)
	m.ClientID = types.StringValue(app.ClientID)
	m.ClientSecret = types.StringValue(app.ClientSecret)
}
