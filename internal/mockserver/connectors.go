package mockserver

import (
	"github.com/apono-io/apono-sdk-go"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func SetupMockHttpServerConnectorV2Api(existingConnectors []apono.Connector) {
	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/connectors", func(req *http.Request) (*http.Response, error) {

		resp, err := httpmock.NewJsonResponse(200, existingConnectors)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}
