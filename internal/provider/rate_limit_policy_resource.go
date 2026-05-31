package provider

import (
	"context"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type rateLimitPolicyResource struct {
	client *bankport.Client
}

type rateLimitPolicyModel struct {
	ID                types.String `tfsdk:"id"`
	ProductCode       types.String `tfsdk:"product_code"`
	SubjectType       types.String `tfsdk:"subject_type"`
	SubjectID         types.String `tfsdk:"subject_id"`
	RequestsPerMinute types.Int64  `tfsdk:"requests_per_minute"`
	BurstLimit        types.Int64  `tfsdk:"burst_limit"`
	Mode              types.String `tfsdk:"mode"`
}

var _ resource.Resource = (*rateLimitPolicyResource)(nil)
var _ resource.ResourceWithConfigure = (*rateLimitPolicyResource)(nil)
var _ resource.ResourceWithImportState = (*rateLimitPolicyResource)(nil)

func NewRateLimitPolicyResource() resource.Resource {
	return &rateLimitPolicyResource{}
}

func (r *rateLimitPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bankport_rate_limit_policy"
}

func (r *rateLimitPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provisions a BankPort API rate-limit policy.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Computed: true},
			"product_code":        schema.StringAttribute{Required: true},
			"subject_type":        schema.StringAttribute{Required: true},
			"subject_id":          schema.StringAttribute{Required: true},
			"requests_per_minute": schema.Int64Attribute{Required: true},
			"burst_limit":         schema.Int64Attribute{Required: true},
			"mode":                schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *rateLimitPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (r *rateLimitPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rateLimitPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := bankport.RateLimitPolicy{
		ProductCode:       plan.ProductCode.ValueString(),
		SubjectType:       plan.SubjectType.ValueString(),
		SubjectID:         plan.SubjectID.ValueString(),
		RequestsPerMinute: plan.RequestsPerMinute.ValueInt64(),
		BurstLimit:        plan.BurstLimit.ValueInt64(),
		Mode:              stringValueOrDefault(plan.Mode, "enforce"),
	}
	created, err := r.client.CreateRateLimitPolicy(ctx, policy)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to create BankPort rate-limit policy", err)
		return
	}

	plan.applyRateLimitPolicy(created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rateLimitPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rateLimitPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetRateLimitPolicy(ctx, state.ID.ValueString())
	if bankport.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to read BankPort rate-limit policy", err)
		return
	}

	state.applyRateLimitPolicy(policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rateLimitPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rateLimitPolicyModel
	var state rateLimitPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := bankport.RateLimitPolicy{
		ProductCode:       plan.ProductCode.ValueString(),
		SubjectType:       plan.SubjectType.ValueString(),
		SubjectID:         plan.SubjectID.ValueString(),
		RequestsPerMinute: plan.RequestsPerMinute.ValueInt64(),
		BurstLimit:        plan.BurstLimit.ValueInt64(),
		Mode:              stringValueOrDefault(plan.Mode, "enforce"),
	}
	updated, err := r.client.UpdateRateLimitPolicy(ctx, state.ID.ValueString(), policy)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to update BankPort rate-limit policy", err)
		return
	}

	plan.applyRateLimitPolicy(updated)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rateLimitPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state rateLimitPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteRateLimitPolicy(ctx, state.ID.ValueString()); err != nil && !bankport.IsNotFound(err) {
		addAPIError(&resp.Diagnostics, "Unable to delete BankPort rate-limit policy", err)
	}
}

func (r *rateLimitPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *rateLimitPolicyModel) applyRateLimitPolicy(policy bankport.RateLimitPolicy) {
	m.ID = types.StringValue(policy.ID)
	m.ProductCode = types.StringValue(policy.ProductCode)
	m.SubjectType = types.StringValue(policy.SubjectType)
	m.SubjectID = types.StringValue(policy.SubjectID)
	m.RequestsPerMinute = types.Int64Value(policy.RequestsPerMinute)
	m.BurstLimit = types.Int64Value(policy.BurstLimit)
	m.Mode = types.StringValue(policy.Mode)
}
