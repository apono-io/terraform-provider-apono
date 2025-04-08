package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigureClientInvoker sets up the client.Invoker from the provider data.
// It's a common utility function to be used in ResourceWithConfigure  DataSourceWithConfigure implementations.
func ConfigureClientInvoker(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse, target *client.Invoker) {
	if req.ProviderData == nil {
		return
	}

	clientProvider, ok := req.ProviderData.(client.ClientProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource or Datasource Configure Type",
			fmt.Sprintf("Expected client.ClientProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	*target = clientProvider.PublicClient()
}
