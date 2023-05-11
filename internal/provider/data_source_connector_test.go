package provider

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/jarcoal/httpmock"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupMockHttpServerConnectorDataSource()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Valid connector ID
			{
				Config: testAccIntegrationDataSourceConfig("test-connector-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.apono_connector.test", "id", "test-connector-id"),
				),
			},
			// Invalid connector ID
			{
				Config:      testAccIntegrationDataSourceConfig("invalid-connector-id"),
				ExpectError: regexp.MustCompile("No connector matched the search criteria"),
			},
		},
	})
}

func testAccIntegrationDataSourceConfig(connectorId string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

data "apono_connector" "test" {
  id = "%[1]s"
}
`, connectorId)
}

func setupMockHttpServerConnectorDataSource() {
	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/connectors", func(req *http.Request) (*http.Response, error) {
		connectors := []apono.Connector{
			{
				ConnectorId:   "test-connector-id",
				LastConnected: *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
				Status:        "Connected",
			},
		}
		resp, err := httpmock.NewJsonResponse(200, connectors)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}
