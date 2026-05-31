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

type sandboxEnvironmentResource struct {
	client *bankport.Client
}

type sandboxEnvironmentModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Products    types.Set    `tfsdk:"products"`
	Region      types.String `tfsdk:"region"`
	Status      types.String `tfsdk:"status"`
	APIKeyToken types.String `tfsdk:"api_key_token"`
}

var _ resource.Resource = (*sandboxEnvironmentResource)(nil)
var _ resource.ResourceWithConfigure = (*sandboxEnvironmentResource)(nil)
var _ resource.ResourceWithImportState = (*sandboxEnvironmentResource)(nil)

func NewSandboxEnvironmentResource() resource.Resource {
	return &sandboxEnvironmentResource{}
}

func (r *sandboxEnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bankport_sandbox_environment"
}

func (r *sandboxEnvironmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provisions a BankPort sandbox environment for partner integration testing.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"name":          schema.StringAttribute{Required: true},
			"products":      schema.SetAttribute{Required: true, ElementType: types.StringType},
			"region":        schema.StringAttribute{Optional: true, Computed: true},
			"status":        schema.StringAttribute{Optional: true, Computed: true},
			"api_key_token": schema.StringAttribute{Computed: true, Sensitive: true},
		},
	}
}

func (r *sandboxEnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (r *sandboxEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sandboxEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env := bankport.SandboxEnvironment{
		Name:     plan.Name.ValueString(),
		Products: stringSetToSlice(ctx, plan.Products, &resp.Diagnostics),
		Region:   stringValueOrDefault(plan.Region, "us-east-1"),
		Status:   stringValueOrDefault(plan.Status, "ready"),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateSandboxEnvironment(ctx, env)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to create BankPort sandbox environment", err)
		return
	}

	plan.applySandboxEnvironment(ctx, created, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sandboxEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sandboxEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env, err := r.client.GetSandboxEnvironment(ctx, state.ID.ValueString())
	if bankport.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to read BankPort sandbox environment", err)
		return
	}

	state.applySandboxEnvironment(ctx, env, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sandboxEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sandboxEnvironmentModel
	var state sandboxEnvironmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	env := bankport.SandboxEnvironment{
		Name:     plan.Name.ValueString(),
		Products: stringSetToSlice(ctx, plan.Products, &resp.Diagnostics),
		Region:   stringValueOrDefault(plan.Region, "us-east-1"),
		Status:   stringValueOrDefault(plan.Status, "ready"),
	}
	if resp.Diagnostics.HasError() {
		return
	}

	updated, err := r.client.UpdateSandboxEnvironment(ctx, state.ID.ValueString(), env)
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to update BankPort sandbox environment", err)
		return
	}

	plan.applySandboxEnvironment(ctx, updated, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sandboxEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sandboxEnvironmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteSandboxEnvironment(ctx, state.ID.ValueString()); err != nil && !bankport.IsNotFound(err) {
		addAPIError(&resp.Diagnostics, "Unable to delete BankPort sandbox environment", err)
	}
}

func (r *sandboxEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (m *sandboxEnvironmentModel) applySandboxEnvironment(ctx context.Context, env bankport.SandboxEnvironment, diags *diag.Diagnostics) {
	m.ID = types.StringValue(env.ID)
	m.Name = types.StringValue(env.Name)
	m.Products = stringSetValue(ctx, env.Products, diags)
	m.Region = types.StringValue(env.Region)
	m.Status = types.StringValue(env.Status)
	m.APIKeyToken = types.StringValue(env.APIKeyToken)
}
