package mockserver

import (
	"encoding/json"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func SetupMockHttpServerManualWebhookEndpoints(existingManualWebhooks []aponoapi.WebhookManualTriggerTerraformModel) {
	var manualWebhooks = map[string]aponoapi.WebhookManualTriggerTerraformModel{}
	for _, manualWebhook := range existingManualWebhooks {
		manualWebhooks[manualWebhook.Id] = manualWebhook
	}

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/terraform/v1/webhooks/manual/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		manualWebhook, exists := manualWebhooks[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Manual Webhook not found"), nil
		}

		resp, err := httpmock.NewJsonResponse(200, manualWebhook)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/terraform/v1/webhooks/manual", func(req *http.Request) (*http.Response, error) {
		var createReq aponoapi.WebhookManualTriggerUpsertTerraformModel
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		manualWebhook := aponoapi.WebhookManualTriggerTerraformModel{
			Id:                           id.String(),
			Name:                         createReq.Name,
			Active:                       createReq.Active,
			BodyTemplate:                 createReq.BodyTemplate,
			ResponseValidators:           createReq.ResponseValidators,
			TimeoutInSec:                 createReq.TimeoutInSec,
			AuthenticationConfig:         createReq.AuthenticationConfig,
			CustomValidationErrorMessage: createReq.CustomValidationErrorMessage,
		}
		manualWebhooks[manualWebhook.Id] = manualWebhook

		resp, err := httpmock.NewJsonResponse(200, manualWebhook)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPut, `=~^http://api\.apono\.dev/api/terraform/v1/webhooks/manual/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		var updateReq aponoapi.WebhookManualTriggerUpsertTerraformModel
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		manualWebhook, exists := manualWebhooks[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Manual Webhook not found"), nil
		}

		manualWebhook.Name = updateReq.Name
		manualWebhook.Active = updateReq.Active
		manualWebhook.BodyTemplate = updateReq.BodyTemplate
		manualWebhook.ResponseValidators = updateReq.ResponseValidators
		manualWebhook.TimeoutInSec = updateReq.TimeoutInSec
		manualWebhook.AuthenticationConfig = updateReq.AuthenticationConfig
		manualWebhook.CustomValidationErrorMessage = updateReq.CustomValidationErrorMessage

		resp, err := httpmock.NewJsonResponse(200, manualWebhook)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/terraform/v1/webhooks/manual/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		_, exists := manualWebhooks[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Manual Webhook not found"), nil
		}

		delete(manualWebhooks, id)

		return httpmock.NewStringResponse(204, ""), nil
	})
}
