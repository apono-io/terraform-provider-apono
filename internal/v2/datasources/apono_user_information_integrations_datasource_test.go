package datasources_test

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAponoUserInformationIntegrationsDataSource(t *testing.T) {
	dataSourceName := "data.apono_user_information_integrations.test"

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(dataSourceName, "integrations.#"),
	}

	if testcommon.IsTestAccount(t) {
		checks = append(checks,
			resource.TestCheckResourceAttr(dataSourceName, "integrations.0.name", "Jumpcloud IDP"),
			resource.TestCheckResourceAttr(dataSourceName, "integrations.0.category", common.UserInformation),
			resource.TestCheckResourceAttr(dataSourceName, "integrations.0.status", "ACTIVE"),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testcommon.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testprovider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAponoUserInformationIntegrationsDataSourceConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

func testAccAponoUserInformationIntegrationsDataSourceConfig() string {
	return `
data "apono_user_information_integrations" "test" {
  name = "Jumpcloud IDP"
}
`
}
