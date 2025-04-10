package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConfigureResourceClientInvoker sets up the client.Invoker from the provider data for resources.
func ConfigureResourceClientInvoker(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, target *client.Invoker) {
	providerData := req.ProviderData

	if providerData == nil {
		return
	}

	clientProvider, ok := providerData.(client.ClientProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected client.ClientProvider, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return
	}

	*target = clientProvider.PublicClient()
}

// ConfigureDataSourceClientInvoker sets up the client.Invoker from the provider data for data sources.
func ConfigureDataSourceClientInvoker(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, target *client.Invoker) {
	providerData := req.ProviderData

	if providerData == nil {
		return
	}

	clientProvider, ok := providerData.(client.ClientProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf("Expected client.ClientProvider, got: %T. Please report this issue to the provider developers.", providerData),
		)
		return
	}

	*target = clientProvider.PublicClient()
}

// SliceToTfValues converts a slice of strings to a slice of TF string values
func SliceToTfValues(values []string) []attr.Value {
	result := make([]attr.Value, 0, len(values))
	for _, v := range values {
		result = append(result, types.StringValue(v))
	}
	return result
}
