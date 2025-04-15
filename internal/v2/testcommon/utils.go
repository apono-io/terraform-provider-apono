package testcommon

import (
	"net/http"
	"os"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/provider"
	v2client "github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"apono": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccPreCheck validates the required environment variables are set.
func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("APONO_PERSONAL_TOKEN"); v == "" {
		t.Fatal("APONO_PERSONAL_TOKEN must be set for acceptance tests")
	}
}

// GetTestClient creates a new (real) Apono API client for acceptance testing.
func GetTestClient(t *testing.T) *v2client.Client {
	endpoint := os.Getenv("APONO_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.apono.io"
	}

	transport := &v2client.DebugTransport{
		Transport: &v2client.UserAgentTransport{
			UserAgent: "terraform-provider-apono/test",
			Transport: http.DefaultTransport,
		},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	token := os.Getenv("APONO_PERSONAL_TOKEN")
	require.NotEmpty(t, token, "APONO_PERSONAL_TOKEN must be set for acceptance tests")

	securitySource := v2client.NewTokenSecuritySource(token)

	client, err := v2client.NewClient(
		endpoint,
		securitySource,
		v2client.WithClient(httpClient),
	)
	require.NoError(t, err, "Failed to create client")

	return client
}

func PrefixedName(prefix, name string) string {
	return prefix + "-" + name
}
