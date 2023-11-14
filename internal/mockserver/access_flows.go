package mockserver

import (
	"encoding/json"
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
	"strings"
	"time"
)

func SetupMockHttpServerAccessFlowV1Endpoints(existingAccessFlows []apono.AccessFlowV1) {
	var accessFlows = map[string]apono.AccessFlowV1{}
	for _, accessFlow := range existingAccessFlows {
		accessFlows[accessFlow.Id] = accessFlow
	}

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/v1/access-flows", func(req *http.Request) (*http.Response, error) {
		var createReq apono.UpsertAccessFlowV1
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessFlow := apono.AccessFlowV1{
			Id:                 id.String(),
			Name:               createReq.Name,
			Active:             createReq.Active,
			Trigger:            createReq.Trigger,
			Grantees:           createReq.Grantees,
			IntegrationTargets: createReq.IntegrationTargets,
			BundleTargets:      createReq.BundleTargets,
			Approvers:          createReq.Approvers,
			RevokeAfterInSec:   createReq.RevokeAfterInSec,
			Settings:           createReq.Settings,
			CreatedDate:        apono.Instant{Time: time.Now()},
		}
		accessFlows[accessFlow.Id] = accessFlow

		fixedJsonResponse, err := fixCreateDateOnJsonResponse(&accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		resp, err := httpmock.NewJsonResponse(200, fixedJsonResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		accessFlow, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		fixedJsonResponse, err := fixCreateDateOnJsonResponse(&accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		resp, err := httpmock.NewJsonResponse(200, fixedJsonResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
	httpmock.RegisterResponder(http.MethodPatch, `=~^http://api\.apono\.dev/api/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq apono.UpdateAccessFlowV1
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		accessFlow, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		if updateReq.Name.IsSet() {
			accessFlow.Name = *updateReq.Name.Get()
		}
		if updateReq.Active.IsSet() {
			accessFlow.Active = *updateReq.Active.Get()
		}
		if updateReq.Trigger.IsSet() {
			accessFlow.Trigger = apono.AccessFlowTriggerV1{
				Type:      updateReq.Trigger.Get().Type,
				Timeframe: updateReq.Trigger.Get().Timeframe,
			}
		}
		if updateReq.HasGrantees() {
			accessFlow.Grantees = updateReq.Grantees
		}
		if updateReq.HasIntegrationTargets() {
			accessFlow.IntegrationTargets = updateReq.IntegrationTargets
		}
		if updateReq.HasBundleTargets() {
			accessFlow.BundleTargets = updateReq.BundleTargets
		}
		if updateReq.HasApprovers() {
			accessFlow.Approvers = updateReq.Approvers
		}
		if updateReq.RevokeAfterInSec.IsSet() {
			accessFlow.RevokeAfterInSec = *updateReq.RevokeAfterInSec.Get()
		}
		if updateReq.Settings.IsSet() {
			accessFlow.Settings = *apono.NewNullableAccessFlowV1Settings(updateReq.Settings.Get())
		}

		fixedJsonResponse, err := fixCreateDateOnJsonResponse(&accessFlow)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		resp, err := httpmock.NewJsonResponse(200, fixedJsonResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessFlows[accessFlow.Id] = accessFlow

		return resp, nil
	})
	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/v1/access-flows/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		_, exists := accessFlows[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Flow not found"), nil
		}

		delete(accessFlows, id)

		messageResponse := apono.MessageResponse{
			Message: "Deleted access flow",
		}

		resp, err := httpmock.NewJsonResponse(200, messageResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func fixCreateDateOnJsonResponse(accessFlow *apono.AccessFlowV1) (map[string]interface{}, error) {
	cleanJson, err := json.Marshal(accessFlow)
	if err != nil {
		return nil, err
	}

	var jsonFields map[string]interface{}
	err = json.Unmarshal(cleanJson, &jsonFields)
	if err != nil {
		return nil, err
	}

	timeAsInstantString, _ := accessFlow.CreatedDate.MarshalJSON()
	jsonFields["created_date"] = strings.Replace(string(timeAsInstantString), "\"", "", -1)

	return jsonFields, nil
}
