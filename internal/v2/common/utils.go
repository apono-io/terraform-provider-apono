package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
