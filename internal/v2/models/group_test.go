package models

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestGroupToDataModel(t *testing.T) {
	tests := []struct {
		name     string
		input    *client.GroupV1
		expected GroupDataModel
	}{
		{
			name: "full group data",
			input: &client.GroupV1{
				ID:                    "group-123",
				Name:                  "Test Group",
				SourceIntegrationID:   client.OptNilString{Value: "integration-123", Set: true},
				SourceIntegrationName: client.OptNilString{Value: "Test Integration", Set: true},
			},
			expected: GroupDataModel{
				ID:                    types.StringValue("group-123"),
				Name:                  types.StringValue("Test Group"),
				SourceIntegrationID:   types.StringValue("integration-123"),
				SourceIntegrationName: types.StringValue("Test Integration"),
			},
		},
		{
			name: "group without source integration",
			input: &client.GroupV1{
				ID:   "group-456",
				Name: "Another Group",
			},
			expected: GroupDataModel{
				ID:   types.StringValue("group-456"),
				Name: types.StringValue("Another Group"),
			},
		},
		{
			name: "group with only source integration ID",
			input: &client.GroupV1{
				ID:                  "group-789",
				Name:                "Third Group",
				SourceIntegrationID: client.OptNilString{Value: "integration-789", Set: true},
			},
			expected: GroupDataModel{
				ID:                  types.StringValue("group-789"),
				Name:                types.StringValue("Third Group"),
				SourceIntegrationID: types.StringValue("integration-789"),
			},
		},
		{
			name: "group with only source integration name",
			input: &client.GroupV1{
				ID:                    "group-abc",
				Name:                  "Fourth Group",
				SourceIntegrationName: client.OptNilString{Value: "Another Integration", Set: true},
			},
			expected: GroupDataModel{
				ID:                    types.StringValue("group-abc"),
				Name:                  types.StringValue("Fourth Group"),
				SourceIntegrationName: types.StringValue("Another Integration"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GroupToDataModel(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
