package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/allanflavio/terraport-go-terraform-provider/internal/bankport"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func clientFromProviderData(data any, diags *diag.Diagnostics) *bankport.Client {
	if data == nil {
		return nil
	}
	client, ok := data.(*bankport.Client)
	if !ok {
		diags.AddError("Unexpected provider data", fmt.Sprintf("Expected *bankport.Client, got %T", data))
		return nil
	}
	return client
}

func stringSetToSlice(ctx context.Context, value types.Set, diags *diag.Diagnostics) []string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var items []string
	diags.Append(value.ElementsAs(ctx, &items, false)...)
	sort.Strings(items)
	return items
}

func stringSetValue(ctx context.Context, values []string, diags *diag.Diagnostics) types.Set {
	sort.Strings(values)
	set, setDiags := types.SetValueFrom(ctx, types.StringType, values)
	diags.Append(setDiags...)
	return set
}

func stringValueOrDefault(value types.String, fallback string) string {
	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		return fallback
	}
	return value.ValueString()
}

func boolValueOrDefault(value types.Bool, fallback bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}
	return value.ValueBool()
}

func addAPIError(diags *diag.Diagnostics, summary string, err error) {
	diags.AddError(summary, bankport.Redact(err.Error()))
}
