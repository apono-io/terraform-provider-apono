package common

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ConfigureClientInvoker sets up the client.Invoker from the provider data.
// Supports both Resource and Datasource Configure implementations.
func ConfigureClientInvoker(ctx context.Context, req interface{}, resp interface{}, target *client.Invoker) {
	var providerData any

	switch v := req.(type) {
	case resource.ConfigureRequest:
		providerData = v.ProviderData
	case datasource.ConfigureRequest:
		providerData = v.ProviderData
	default:
		panic(fmt.Sprintf("Unsupported request type: %T", req))
	}

	if providerData == nil {
		return
	}

	clientProvider, ok := providerData.(client.ClientProvider)
	if !ok {
		switch v := resp.(type) {
		case *resource.ConfigureResponse:
			v.Diagnostics.AddError(
				"Unexpected Configure Type",
				fmt.Sprintf("Expected client.ClientProvider, got: %T. Please report this issue to the provider developers.", providerData),
			)
		case *datasource.ConfigureResponse:
			v.Diagnostics.AddError(
				"Unexpected Configure Type",
				fmt.Sprintf("Expected client.ClientProvider, got: %T. Please report this issue to the provider developers.", providerData),
			)
		default:
			panic(fmt.Sprintf("Unsupported response type: %T", resp))
		}
		return
	}

	*target = clientProvider.PublicClient()
}
