package provider

import (
	"encoding/json"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIntegrationResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupMockHttpServerIntegrationResource()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccIntegrationResourceConfig("integration-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("apono_integration.test", "id"),
					resource.TestCheckResourceAttr("apono_integration.test", "name", "integration-name"),
					resource.TestCheckResourceAttr("apono_integration.test", "type", "postgresql"),
					resource.TestCheckResourceAttr("apono_integration.test", "aws_secret.region", "us-east-1"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "apono_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIntegrationResourceConfig("updated-name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apono_integration.test", "name", "updated-name"),
					resource.TestCheckResourceAttr("apono_integration.test", "type", "postgresql"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIntegrationResourceConfig(integrationName string) string {
	return fmt.Sprintf(`
provider apono {
  endpoint = "http://api.apono.dev"
  personal_token = "1234567890abcdefg"
}

resource "apono_integration" "test" {
  name = "%[1]s"
  type = "postgresql"
  connector_id = "000-1111-222222-33333-444444"
  metadata = {
    hostname = "my-postgres-rds.aaabbbsss111.us-east-1.rds.amazonaws.com"
    port = "5432"
    dbname = "postgres"
  }
  aws_secret = {
    region = "us-east-1"
    secret_id = "my-secret"
  }
}
`, integrationName)
}

func setupMockHttpServerIntegrationResource() {
	var integrations = map[string]*apono.Integration{}
	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/v2/integrations", func(req *http.Request) (*http.Response, error) {
		var createReq apono.CreateIntegration
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		integration := apono.Integration{
			Id:            id.String(),
			Name:          createReq.Name,
			Type:          createReq.Type,
			ProvisionerId: createReq.ProvisionerId,
			Status:        apono.INTEGRATIONSTATUS_ACTIVE,
			Metadata:      createReq.Metadata,
			SecretConfig:  createReq.SecretConfig,
		}
		integrations[integration.Id] = &integration

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		integration, exists := integrations[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Integration not found"), nil
		}

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPut, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq apono.UpdateIntegration
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		integration, exists := integrations[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Integration not found"), nil
		}

		integration.Name = updateReq.Name
		integration.ProvisionerId = updateReq.ProvisionerId
		integration.Metadata = updateReq.Metadata
		integration.SecretConfig = updateReq.SecretConfig

		resp, err := httpmock.NewJsonResponse(200, integration)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/v2/integrations/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		delete(integrations, id)

		messageResponse := apono.MessageResponse{
			Message: "Deleted integration",
		}

		resp, err := httpmock.NewJsonResponse(200, messageResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v2/integrations-catalog/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		configType := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		switch configType {
		case "postgresql":
			config := apono.IntegrationConfig{
				Name:        "PostgreSQL",
				Type:        "postgresql",
				Description: "An open-source relational database management system emphasizing extensibility and SQL compliance.",
				Params: []apono.IntegrationConfigParam{
					{
						Id:    "hostname",
						Label: "Hostname",
					},
					{
						Id:      "port",
						Label:   "Port",
						Default: "5432",
					},
					{
						Id:      "dbname",
						Label:   "Database Name",
						Default: "postgres",
					},
				},
				RequiresSecret: true,
				SupportedSecretTypes: []string{
					"AWS",
					"GCP",
					"KUBERNETES",
				},
			}

			resp, err := httpmock.NewJsonResponse(200, config)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}

			return resp, nil
		default:
			return httpmock.NewStringResponse(400, "Unsupported config type"), nil
		}
	})
}
