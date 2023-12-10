package mockserver

import (
	"encoding/json"
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"golang.org/x/exp/maps"
	"net/http"
)

func SetupMockHttpServerAccessBundleV1Endpoints(existingAccessBundles []apono.AccessBundleV1) {
	var accessBundles = map[string]apono.AccessBundleV1{}
	for _, accessBundle := range existingAccessBundles {
		accessBundles[accessBundle.Id] = accessBundle
	}

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v1/access-bundles/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		accessBundle, exists := accessBundles[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Bundle not found"), nil
		}

		resp, err := httpmock.NewJsonResponse(200, accessBundle)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v1/access-bundles", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, apono.PaginatedResponseAccessBundleV1Model{
			Data: maps.Values(accessBundles),
			Pagination: apono.PaginationInfo{
				Total:  int32(len(accessBundles)),
				Limit:  int32(len(accessBundles)),
				Offset: 0,
			},
		})
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPost, "http://api.apono.dev/api/v1/access-bundles", func(req *http.Request) (*http.Response, error) {
		var createReq apono.UpsertAccessBundleV1
		if err := json.NewDecoder(req.Body).Decode(&createReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		id, err := uuid.NewUUID()
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessBundle := apono.AccessBundleV1{
			Id:                 id.String(),
			Name:               createReq.Name,
			IntegrationTargets: createReq.IntegrationTargets,
		}
		accessBundles[accessBundle.Id] = accessBundle

		resp, err := httpmock.NewJsonResponse(200, accessBundle)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodPatch, `=~^http://api\.apono\.dev/api/v1/access-bundles/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		var updateReq apono.UpdateAccessBundleV1
		if err := json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
			return httpmock.NewStringResponse(400, err.Error()), nil
		}

		accessBundle, exists := accessBundles[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Bundle not found"), nil
		}

		if updateReq.Name.IsSet() {
			accessBundle.Name = *updateReq.Name.Get()
		}

		if updateReq.HasIntegrationTargets() {
			accessBundle.IntegrationTargets = updateReq.IntegrationTargets
		}

		resp, err := httpmock.NewJsonResponse(200, accessBundle)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		accessBundles[accessBundle.Id] = accessBundle

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodDelete, `=~^http://api\.apono\.dev/api/v1/access-bundles/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch

		_, exists := accessBundles[id]
		if !exists {
			return httpmock.NewStringResponse(404, "Access Bundle not found"), nil
		}

		delete(accessBundles, id)

		messageResponse := apono.MessageResponse{
			Message: "Deleted access bundle",
		}

		resp, err := httpmock.NewJsonResponse(200, messageResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func CreateMockAccessBundles() []apono.AccessBundleV1 {
	return []apono.AccessBundleV1{
		{
			Id:   uuid.NewString(),
			Name: "DB DEV",
			IntegrationTargets: []apono.AccessTargetIntegrationV1{
				{
					IntegrationId:       "1",
					ResourceType:        "mysql-db",
					ResourceTagIncludes: nil,
					ResourceTagExcludes: []apono.TagV1{{Name: "env", Value: "prod"}},
					Permissions:         []string{"ReadOnly", "ReadWrite", "Admin"},
				},
				{
					IntegrationId:       "2",
					ResourceType:        "postgresql-db",
					ResourceTagIncludes: []apono.TagV1{{Name: "env", Value: "dev"}},
					ResourceTagExcludes: nil,
					Permissions:         []string{"ReadOnly", "ReadWrite", "Admin"},
				},
			},
		},
		{
			Id:   uuid.NewString(),
			Name: "DB PROD",
			IntegrationTargets: []apono.AccessTargetIntegrationV1{
				{
					IntegrationId:       "3",
					ResourceType:        "mysql-db",
					ResourceTagIncludes: nil,
					ResourceTagExcludes: []apono.TagV1{{Name: "env", Value: "dev"}},
					Permissions:         []string{"ReadOnly", "ReadWrite", "Admin"},
				},
				{
					IntegrationId:       "4",
					ResourceType:        "postgresql-db",
					ResourceTagIncludes: []apono.TagV1{{Name: "env", Value: "prod"}},
					ResourceTagExcludes: nil,
					Permissions:         []string{"ReadOnly", "ReadWrite", "Admin"},
				},
			},
		},
	}
}
