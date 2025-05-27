package testcommon

import (
	"net/http"
	"os"
	"testing"

	v2client "github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

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

func IsTestAccount(t *testing.T) bool {
	return os.Getenv("IS_TEST_ACCOUNT") != ""
}

func CreateTestStringSet(t *testing.T, values []string) types.Set {
	result, diags := types.SetValueFrom(t.Context(), types.StringType, values)
	require.False(t, diags.HasError())
	return result
}

func CreateTestStringList(t *testing.T, values []string) types.List {
	result, diags := types.ListValueFrom(t.Context(), types.StringType, values)
	require.False(t, diags.HasError())
	return result
}
