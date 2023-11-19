package provider

import (
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/mockserver"
	"github.com/jarcoal/httpmock"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConnectorDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	connectors := []apono.Connector{
		{
			ConnectorId:   "test-connector-id",
			LastConnected: *apono.NewNullableInstant(&apono.Instant{Time: time.Now()}),
			Status:        "Connected",
		},
	}
	mockserver.SetupMockHttpServerConnectorV2Api(connectors)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Valid connector ID
			{
				Config: testAccConnectorDataSourceConfig("test-connector-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.apono_connector.test", "id", "test-connector-id"),
				),
			},
			// Invalid connector ID
			{
				Config:      testAccConnectorDataSourceConfig("invalid-connector-id"),
				ExpectError: regexp.MustCompile("No connector matched the search criteria"),
			},
		},
	})
}

func testAccConnectorDataSourceConfig(connectorId string) string {
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
