package testcommon

import (
	"os"
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"apono": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccPreCheck validates the required environment variables are set
func TestAccPreCheck(t *testing.T) {
	// Check for required environment variables
	if v := os.Getenv("APONO_ENDPOINT"); v == "" {
		t.Fatal("APONO_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("APONO_PERSONAL_TOKEN"); v == "" {
		t.Fatal("APONO_PERSONAL_TOKEN must be set for acceptance tests")
	}
}
