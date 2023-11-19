package mockserver

import (
	"github.com/apono-io/apono-sdk-go"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func SetupMockHttpServerUsersV2Endpoints(existingUsers []apono.UserModel) {
	var usersById = map[string]apono.UserModel{}
	for _, user := range existingUsers {
		usersById[user.Id] = user
	}
	var usersByEmail = map[string]apono.UserModel{}
	for _, user := range existingUsers {
		usersByEmail[user.Email] = user
	}

	httpmock.RegisterResponder(http.MethodGet, `=~^http://api\.apono\.dev/api/v2/users/([^/]+)\z`, func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1) // 1=first regexp submatch
		userById, existsById := usersById[id]
		userByEmail, existsByEmail := usersByEmail[id]
		if !existsById && !existsByEmail {
			return httpmock.NewStringResponse(404, "User not found"), nil
		}

		var user apono.UserModel
		if existsById {
			user = userById
		} else {
			user = userByEmail
		}

		resp, err := httpmock.NewJsonResponse(200, user)
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})

	httpmock.RegisterResponder(http.MethodGet, "http://api.apono.dev/api/v2/users", func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, apono.PaginatedResponseUserModel{
			Data: existingUsers,
			Pagination: apono.PaginationInfo{
				Total:  int32(len(existingUsers)),
				Limit:  int32(len(existingUsers)),
				Offset: 0,
			},
		})
		if err != nil {
			return httpmock.NewStringResponse(500, err.Error()), nil
		}

		return resp, nil
	})
}

func CreateMockUsers() []apono.UserModel {
	return []apono.UserModel{
		{
			Id:        uuid.NewString(),
			Email:     "test1@example.com",
			FirstName: "Test",
			LastName:  "User 1",
			Active:    true,
		},
		{
			Id:        uuid.NewString(),
			Email:     "test2@example.com",
			FirstName: "Test",
			LastName:  "User 2",
			Active:    true,
		},
	}
}
