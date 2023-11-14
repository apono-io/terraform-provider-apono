package mockserver

import (
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func SetupMockHttpServerIdentitiesV2Endpoints(existingIdentities []apono.IdentityModel2) {
	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/identities", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, apono.PaginatedResponseIdentityModelV2{
			Data: existingIdentities,
			Pagination: apono.PaginationInfo{
				Total:  int32(len(existingIdentities)),
				Limit:  int32(len(existingIdentities)),
				Offset: 0,
			},
		})
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func CreateMockIdentities() []apono.IdentityModel2 {
	return []apono.IdentityModel2{
		{
			Id:   uuid.NewString(),
			Name: "Test Group 1",
			Type: "GROUP",
		},
		{
			Id:   uuid.NewString(),
			Name: "Test Group 2",
			Type: "GROUP",
		},
		{
			Id:   uuid.NewString(),
			Name: "Test User 1",
			Type: "USER",
		},
		{
			Id:   uuid.NewString(),
			Name: "Test User 2",
			Type: "USER",
		},
		{
			Id:   uuid.NewString(),
			Name: "Manager",
			Type: "CONTEXT_ATTRIBUTE",
		},
		{
			Id:   uuid.NewString(),
			Name: "Shift",
			Type: "CONTEXT_ATTRIBUTE",
		},
	}
}
