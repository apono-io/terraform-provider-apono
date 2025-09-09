package testcommon

import (
	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func GenerateAccessFlowResponse() *client.AccessFlowV2 {
	response := client.AccessFlowV2{
		ID:      "flow-123",
		Name:    "postgresql_prod",
		Active:  true,
		Trigger: "SELF_SERVE",
		Settings: client.AccessFlowSettingsV2{
			JustificationRequired:         true,
			RequireApproverReason:         false,
			RequestorCannotApproveHimself: false,
			RequireMfa:                    false,
			Labels:                        []string{"DB", "PROD", "TERRAFORM"},
		},
	}

	response.GrantDurationInMin.SetTo(int32(60))

	timeframe := client.AccessFlowTimeframeV2{
		StartTime:  "10:00",
		EndTime:    "23:59",
		TimeZone:   "Asia/Jerusalem",
		DaysOfWeek: []client.DayOfWeekV2{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"},
	}
	response.Timeframe.SetTo(timeframe)

	response.Requestors = client.RequestorsV2{
		LogicalOperator: "OR",
		Conditions: []client.ConditionV2{
			{
				Type: "user",
			},
		},
	}
	response.Requestors.Conditions[0].SourceIntegrationName.SetTo("Okta Directory")
	response.Requestors.Conditions[0].MatchOperator.SetTo("is")
	response.Requestors.Conditions[0].Values.SetTo([]string{"person@example.com", "person_two@example.com"})

	requestFor := client.RequestForV2{
		RequestScopes: []string{"self", "others"},
	}
	grantees := client.GranteesV2{
		LogicalOperator: "OR",
		Conditions: []client.ConditionV2{
			{
				Type: "user",
			},
		},
	}
	grantees.Conditions[0].SourceIntegrationName.SetTo("Google Oauth")
	grantees.Conditions[0].MatchOperator.SetTo("is")
	grantees.Conditions[0].Values.SetTo([]string{"user1@example.com"})
	requestFor.Grantees.SetTo(grantees)
	response.RequestFor.SetTo(requestFor)

	bundleTarget := client.AccessTargetV2{}
	bundleData := client.BundleAccessTargetV2{
		BundleID:   "bundle-123",
		BundleName: "PROD ENV",
	}
	bundleTarget.Bundle.SetTo(bundleData)

	integrationTarget := client.AccessTargetV2{}
	integrationData := client.IntegrationAccessTargetV2{
		IntegrationID:   "integration-123",
		IntegrationName: "postgresql",
		ResourceType:    "database",
		Permissions:     []string{"read", "write"},
	}

	resourceScope := client.ResourcesScopeIntegrationAccessTargetV2{
		ScopeMode: "include_resources",
		Type:      "NAME",
		Values:    []string{"db1", "db2"},
	}
	integrationData.ResourcesScopes.SetTo([]client.ResourcesScopeIntegrationAccessTargetV2{resourceScope})
	integrationTarget.Integration.SetTo(integrationData)

	accessScopeTarget := client.AccessTargetV2{}
	accessScopeData := client.AccessScopeAccessTargetV2{
		AccessScopeID:   "scope-123",
		AccessScopeName: "Test Scope",
	}
	accessScopeTarget.AccessScope.SetTo(accessScopeData)

	response.AccessTargets = []client.AccessTargetV2{bundleTarget, integrationTarget, accessScopeTarget}

	approverPolicy := client.ApproverPolicyV2{
		ApprovalMode: "ANY_OF",
		ApproverGroups: []client.ApproverGroupV2{
			{
				LogicalOperator: "OR",
				Approvers: []client.ConditionV2{
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
