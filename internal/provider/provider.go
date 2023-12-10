package provider

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/url"
	"os"
)

// Ensure AponoProvider satisfies various provider interfaces.
var _ provider.Provider = &AponoProvider{}

// AponoProvider defines the provider implementation.
type AponoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	client  *apono.APIClient
}

// AponoProviderConfig describes the provider data model.
type AponoProviderConfig struct {
	Endpoint      types.String `tfsdk:"endpoint"`
	PersonalToken types.String `tfsdk:"personal_token"`
}

func (p *AponoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "apono"
	resp.Version = p.version
}

func (p *AponoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "Override API endpoint",
				Optional:    true,
			},
			"personal_token": schema.StringAttribute{
				Description: "Personal API token",
				Required:    true,
			},
		},
	}
}

func (p *AponoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Start configuring Apono provider")

	// Check environment variables
	endpoint := os.Getenv("APONO_ENDPOINT")
	personalToken := os.Getenv("APONO_PERSONAL_TOKEN")

	var config AponoProviderConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.ValueString() != "" {
		endpoint = config.Endpoint.ValueString()
	}

	if config.PersonalToken.ValueString() != "" {
		personalToken = config.PersonalToken.ValueString()
	}

	if endpoint == "" {
		endpoint = "https://api.apono.io"
	}

	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Endpoint Configuration",
			fmt.Sprintf("Failed to parse endpoint %s: %s", endpoint, err.Error()),
		)
	}

	if personalToken == "" {
		resp.Diagnostics.AddError(
			"Missing Personal API Token Configuration",
			"While configuring the provider, the Personal API token was not found in "+
				"the APONO_PERSONAL_TOKEN environment variable or provider "+
				"configuration block personal_token attribute.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	cfg := apono.NewConfiguration()
	cfg.Scheme = endpointUrl.Scheme
	cfg.Host = endpointUrl.Host
	cfg.UserAgent = fmt.Sprintf("terraform-provider-apono/%s", p.version)
	cfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", personalToken))

	p.client = apono.NewAPIClient(cfg)
	tflog.Debug(ctx, "Provider configuration", map[string]interface{}{
		"provider": fmt.Sprintf("%+v", p),
	})

	resp.DataSourceData = p
	resp.ResourceData = p
}

func (p *AponoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewIntegrationResource,
		NewAccessFlowResource,
		NewAccessBundleResource,
	}
}

func (p *AponoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConnectorDataSource,
		NewIntegrationsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AponoProvider{
			version: version,
		}
	}
}

// toProvider can be used to cast a generic provider.Provider reference to this specific provider.
// This is ideally used in DataSourceType.NewDataSource and ResourceType.NewResource calls.
func toProvider(in any) (*AponoProvider, diag.Diagnostics) {
	if in == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	p, ok := in.(*AponoProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. "+
				"This is always a bug in the provider code and should be reported to the provider developers.", in,
			),
		)
		return nil, diags
	}

	return p, diags
}
