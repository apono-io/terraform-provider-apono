package mockserver

import (
	"encoding/json"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
	"time"
)

func SetupMockHttpServerAccessFlowV1Endpoints(existingAccessFlows []aponoapi.AccessFlowTerraformV1) {
	var accessFlows = map[string]aponoapi.AccessFlowTerraformV1{}
	for _, accessFlow := range existingAccessFlows {
		accessFlows[accessFlow.Id] = accessFlow
	}

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/terraform/v1/access-flows", func(req *http.Request) (*http.Response, error) {
		var createReq aponoapi.UpsertAccessFlowTerraformV1
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessFlow := aponoapi.AccessFlowTerraformV1{
			Id:                 id.String(),
			Name:               createReq.Name,
			Active:             createReq.Active,
			Trigger:            createReq.Trigger,
			Grantees:           createReq.Grantees,
			GranteeFilterGroup: createReq.GranteeFilterGroup,
			IntegrationTargets: createReq.IntegrationTargets,
			BundleTargets:      createReq.BundleTargets,
			Approvers:          createReq.Approvers,
			RevokeAfterInSec:   createReq.RevokeAfterInSec,
			Settings:           createReq.Settings,
			CreatedDate:        getTimeAsInstantFloat(),
		}
		accessFlows[accessFlow.Id] = accessFlow

		resp, err := httpmock.NewJsonResponse(200, accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/terraform/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		accessFlow, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		resp, err := httpmock.NewJsonResponse(200, accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPut, `=~^http://api\.apono\.dev/api/terraform/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq aponoapi.UpsertAccessFlowTerraformV1
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		accessFlow, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		accessFlow.Name = updateReq.Name
		accessFlow.Active = updateReq.Active
		accessFlow.Trigger = aponoapi.AccessFlowTriggerTerraformV1{
			Type:      updateReq.Trigger.Type,
			Timeframe: updateReq.Trigger.Timeframe,
		}
		accessFlow.Grantees = updateReq.Grantees
		accessFlow.GranteeFilterGroup = updateReq.GranteeFilterGroup
		accessFlow.IntegrationTargets = updateReq.IntegrationTargets
		accessFlow.BundleTargets = updateReq.BundleTargets
		accessFlow.Approvers = updateReq.Approvers
		accessFlow.RevokeAfterInSec = updateReq.RevokeAfterInSec
		accessFlow.Settings = updateReq.Settings

		resp, err := httpmock.NewJsonResponse(200, accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessFlows[accessFlow.Id] = accessFlow

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/terraform/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		_, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		delete(accessFlows, id)

		messageResponse := aponoapi.MessageResponse{
			Message: "Deleted access flow",
		}

		resp, err := httpmock.NewJsonResponse(200, messageResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func getTimeAsInstantFloat() float64 {
	now := time.Now()

	milliseconds := float64(now.UnixMilli())

	nanoseconds := float64(now.Nanosecond()) / 1_000_000_000.0

	return milliseconds + nanoseconds
}
