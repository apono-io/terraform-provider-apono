package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func GenerateAccessFlowResponse() *client.AccessFlowPublicV2Model {
	response := client.AccessFlowPublicV2Model{
		ID:      "flow-123",
		Name:    "postgresql_prod",
		Active:  true,
		Trigger: "SELF_SERVE",
		Settings: client.AccessFlowSettingsPublicV2Model{
			JustificationRequired:         true,
			RequireApproverReason:         false,
			RequestorCannotApproveHimself: false,
			RequireMfa:                    false,
			Labels:                        []string{"DB", "PROD", "TERRAFORM"},
		},
	}

	response.GrantDurationInMin.SetTo(int32(60))

	timeframe := client.AccessFlowTimeframePublicV2Model{
		StartTime:  "10:00",
		EndTime:    "23:59",
		TimeZone:   "Asia/Jerusalem",
		DaysOfWeek: []client.DayOfWeekPublicV2Model{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"},
	}
	response.Timeframe.SetTo(timeframe)

	response.Requestors = client.RequestorsPublicV2Model{
		LogicalOperator: "OR",
		Conditions: []client.ConditionPublicV2Model{
			{
				Type: "user",
			},
		},
	}
	response.Requestors.Conditions[0].SourceIntegrationName.SetTo("Okta Directory")
	response.Requestors.Conditions[0].MatchOperator.SetTo("is")
	response.Requestors.Conditions[0].Values.SetTo([]string{"person@example.com", "person_two@example.com"})

	bundleTarget := client.AccessTargetPublicV2Model{}
	bundleData := client.BundleAccessTargetPublicV2Model{
		BundleID:   "bundle-123",
		BundleName: "PROD ENV",
	}
	bundleTarget.Bundle.SetTo(bundleData)

	integrationTarget := client.AccessTargetPublicV2Model{}
	integrationData := client.IntegrationAccessTargetPublicV2Model{
		IntegrationID:   "integration-123",
		IntegrationName: "postgresql",
		ResourceType:    "database",
		Permissions:     []string{"read", "write"},
	}

	resourceScope := client.ResourcesScopeIntegrationAccessTargetPublicV2Model{
		ScopeMode: "include_resources",
		Type:      "NAME",
		Values:    []string{"db1", "db2"},
	}
	integrationData.ResourcesScopes.SetTo([]client.ResourcesScopeIntegrationAccessTargetPublicV2Model{resourceScope})
	integrationTarget.Integration.SetTo(integrationData)

	accessScopeTarget := client.AccessTargetPublicV2Model{}
	accessScopeData := client.AccessScopeAccessTargetPublicV2Model{
		AccessScopeID:   "scope-123",
		AccessScopeName: "Test Scope",
	}
	accessScopeTarget.AccessScope.SetTo(accessScopeData)

	response.AccessTargets = []client.AccessTargetPublicV2Model{bundleTarget, integrationTarget, accessScopeTarget}

	approverPolicy := client.ApproverPolicyPublicV2Model{
		ApprovalMode: "ANY_OF",
		ApproverGroups: []client.ApproverGroupPublicV2Model{
			{
				LogicalOperator: "OR",
				Approvers: []client.ConditionPublicV2Model{
					{
						Type: "user",
					},
				},
			},
		},
	}
	approverPolicy.ApproverGroups[0].Approvers[0].SourceIntegrationName.SetTo("Okta Directory")
	approverPolicy.ApproverGroups[0].Approvers[0].MatchOperator.SetTo("is")
	approverPolicy.ApproverGroups[0].Approvers[0].Values.SetTo([]string{"person@example.com", "person_two@example.com"})
	response.ApproverPolicy.SetTo(approverPolicy)

	return &response
}
