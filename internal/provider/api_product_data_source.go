package provider

import (
	"context"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type apiProductDataSource struct {
	client *bankport.Client
}

type apiProductModel struct {
	ProductCode  types.String `tfsdk:"product_code"`
	Name         types.String `tfsdk:"name"`
	Category     types.String `tfsdk:"category"`
	Beta         types.Bool   `tfsdk:"beta"`
	Regions      types.Set    `tfsdk:"regions"`
	Capabilities types.Set    `tfsdk:"capabilities"`
	DocsURL      types.String `tfsdk:"docs_url"`
}

var _ datasource.DataSource = (*apiProductDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*apiProductDataSource)(nil)

func NewAPIProductDataSource() datasource.DataSource {
	return &apiProductDataSource{}
}

func (d *apiProductDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bankport_api_product"
}

func (d *apiProductDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads BankPort platform API product metadata.",
		Attributes: map[string]schema.Attribute{
			"product_code": schema.StringAttribute{Required: true},
			"name":         schema.StringAttribute{Computed: true},
			"category":     schema.StringAttribute{Computed: true},
			"beta":         schema.BoolAttribute{Computed: true},
			"regions":      schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"capabilities": schema.SetAttribute{Computed: true, ElementType: types.StringType},
			"docs_url":     schema.StringAttribute{Computed: true},
		},
	}
}

func (d *apiProductDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = clientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

func (d *apiProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config apiProductModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	product, err := d.client.GetAPIProduct(ctx, config.ProductCode.ValueString())
	if err != nil {
		addAPIError(&resp.Diagnostics, "Unable to read BankPort API product", err)
		return
	}

	config.Name = types.StringValue(product.Name)
	config.Category = types.StringValue(product.Category)
	config.Beta = types.BoolValue(product.Beta)
	config.Regions = stringSetValue(ctx, product.Regions, &resp.Diagnostics)
	config.Capabilities = stringSetValue(ctx, product.Capabilities, &resp.Diagnostics)
	config.DocsURL = types.StringValue(product.DocsURL)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
