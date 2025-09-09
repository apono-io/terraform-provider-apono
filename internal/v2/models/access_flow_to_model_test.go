package models

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessFlowResponseToModel(t *testing.T) {
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
	integrationData.ResourcesScopes = client.NewOptNilResourcesScopeIntegrationAccessTargetV2Array([]client.ResourcesScopeIntegrationAccessTargetV2{resourceScope})
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

	ctx := t.Context()
	model, err := AccessFlowResponseToModel(ctx, response)
	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "flow-123", model.ID.ValueString())
	assert.Equal(t, "postgresql_prod", model.Name.ValueString())
	assert.True(t, model.Active.ValueBool())
	assert.Equal(t, "SELF_SERVE", model.Trigger.ValueString())
	assert.Equal(t, int32(60), model.GrantDurationInMin.ValueInt32())

	require.NotNil(t, model.Timeframe)
	assert.Equal(t, "10:00", model.Timeframe.StartTime.ValueString())
	assert.Equal(t, "23:59", model.Timeframe.EndTime.ValueString())
	assert.Equal(t, "Asia/Jerusalem", model.Timeframe.TimeZone.ValueString())

	var daysOfWeek []string
	diags := model.Timeframe.DaysOfWeek.ElementsAs(ctx, &daysOfWeek, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"}, daysOfWeek)

	require.NotNil(t, model.Settings)
	assert.True(t, model.Settings.JustificationRequired.ValueBool())
	assert.False(t, model.Settings.RequireApproverReason.ValueBool())
	assert.False(t, model.Settings.RequesterCannotApproveSelf.ValueBool())
	assert.False(t, model.Settings.RequireMFA.ValueBool())

	var labels []string
	diags = model.Settings.Labels.ElementsAs(ctx, &labels, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"DB", "PROD", "TERRAFORM"}, labels)

	require.NotNil(t, model.Requestors)
	assert.Equal(t, "OR", model.Requestors.LogicalOperator.ValueString())
	require.Len(t, model.Requestors.Conditions, 1)
	assert.Equal(t, "Okta Directory", model.Requestors.Conditions[0].SourceIntegrationName.ValueString())
	assert.Equal(t, "user", model.Requestors.Conditions[0].Type.ValueString())
	assert.Equal(t, "is", model.Requestors.Conditions[0].MatchOperator.ValueString())

	var values []string
	diags = model.Requestors.Conditions[0].Values.ElementsAs(ctx, &values, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"person@example.com", "person_two@example.com"}, values)

	require.Len(t, model.AccessTargets, 3)

	require.NotNil(t, model.AccessTargets[0].Bundle)
	assert.Equal(t, "PROD ENV", model.AccessTargets[0].Bundle.Name.ValueString())

	require.NotNil(t, model.AccessTargets[1].Integration)
	assert.Equal(t, "postgresql", model.AccessTargets[1].Integration.IntegrationName.ValueString())
	assert.Equal(t, "database", model.AccessTargets[1].Integration.ResourceType.ValueString())

	var permissions []string
	diags = model.AccessTargets[1].Integration.Permissions.ElementsAs(ctx, &permissions, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"read", "write"}, permissions)

	require.Len(t, model.AccessTargets[1].Integration.ResourcesScopes, 1)
	assert.Equal(t, "include_resources", model.AccessTargets[1].Integration.ResourcesScopes[0].ScopeMode.ValueString())
	assert.Equal(t, "NAME", model.AccessTargets[1].Integration.ResourcesScopes[0].Type.ValueString())

	var scopeValues []string
	diags = model.AccessTargets[1].Integration.ResourcesScopes[0].Values.ElementsAs(ctx, &scopeValues, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"db1", "db2"}, scopeValues)

	require.NotNil(t, model.AccessTargets[2].AccessScope)
	assert.Equal(t, "Test Scope", model.AccessTargets[2].AccessScope.Name.ValueString())

	require.NotNil(t, model.ApproverPolicy)
	assert.Equal(t, "ANY_OF", model.ApproverPolicy.ApprovalMode.ValueString())
	require.Len(t, model.ApproverPolicy.ApproverGroups, 1)
	assert.Equal(t, "OR", model.ApproverPolicy.ApproverGroups[0].LogicalOperator.ValueString())
	require.Len(t, model.ApproverPolicy.ApproverGroups[0].Approvers, 1)
	assert.Equal(t, "Okta Directory", model.ApproverPolicy.ApproverGroups[0].Approvers[0].SourceIntegrationName.ValueString())
	assert.Equal(t, "user", model.ApproverPolicy.ApproverGroups[0].Approvers[0].Type.ValueString())
	assert.Equal(t, "is", model.ApproverPolicy.ApproverGroups[0].Approvers[0].MatchOperator.ValueString())

	var approverValues []string
	diags = model.ApproverPolicy.ApproverGroups[0].Approvers[0].Values.ElementsAs(ctx, &approverValues, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"person@example.com", "person_two@example.com"}, approverValues)

	require.NotNil(t, model.RequestFor)
	var requestScopes []string
	diags = model.RequestFor.RequestScopes.ElementsAs(ctx, &requestScopes, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"self", "others"}, requestScopes)

	require.NotNil(t, model.RequestFor.Grantees)
	assert.Equal(t, "OR", model.RequestFor.Grantees.LogicalOperator.ValueString())
	require.Len(t, model.RequestFor.Grantees.Conditions, 1)
	assert.Equal(t, "Google Oauth", model.RequestFor.Grantees.Conditions[0].SourceIntegrationName.ValueString())
	assert.Equal(t, "user", model.RequestFor.Grantees.Conditions[0].Type.ValueString())
	assert.Equal(t, "is", model.RequestFor.Grantees.Conditions[0].MatchOperator.ValueString())

	var granteeValues []string
	diags = model.RequestFor.Grantees.Conditions[0].Values.ElementsAs(ctx, &granteeValues, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"user1@example.com"}, granteeValues)
}

func TestAccessFlowResponseToModelMinimalFields(t *testing.T) {
	response := client.AccessFlowV2{
		ID:      "flow-456",
		Name:    "minimal_flow",
		Active:  false,
		Trigger: "AUTOMATIC",
		Settings: client.AccessFlowSettingsV2{
			JustificationRequired:         false,
			RequireApproverReason:         false,
			RequestorCannotApproveHimself: false,
			RequireMfa:                    false,
			Labels:                        []string{},
		},
		Requestors: client.RequestorsV2{
			LogicalOperator: "AND",
			Conditions: []client.ConditionV2{
				{
					Type: "user",
				},
			},
		},
	}
	response.Requestors.Conditions[0].MatchOperator.SetTo("is")
	response.Requestors.Conditions[0].Values.SetTo([]string{"person@example.com"})

	bundleTarget := client.AccessTargetV2{}
	bundleData := client.BundleAccessTargetV2{
		BundleID:   "bundle-456",
		BundleName: "QA ENV",
	}
	bundleTarget.Bundle.SetTo(bundleData)
	response.AccessTargets = []client.AccessTargetV2{bundleTarget}

	ctx := t.Context()
	model, err := AccessFlowResponseToModel(ctx, response)
	require.NoError(t, err)
	require.NotNil(t, model)

	assert.Equal(t, "flow-456", model.ID.ValueString())
	assert.Equal(t, "minimal_flow", model.Name.ValueString())
	assert.False(t, model.Active.ValueBool())
	assert.Equal(t, "AUTOMATIC", model.Trigger.ValueString())
	assert.True(t, model.GrantDurationInMin.IsNull())
	assert.Nil(t, model.Timeframe)
	assert.Nil(t, model.ApproverPolicy)

	require.NotNil(t, model.Settings)
	assert.False(t, model.Settings.JustificationRequired.ValueBool())
	assert.False(t, model.Settings.RequireApproverReason.ValueBool())
	assert.False(t, model.Settings.RequesterCannotApproveSelf.ValueBool())
	assert.False(t, model.Settings.RequireMFA.ValueBool())
	assert.True(t, model.Settings.Labels.IsNull())

	require.NotNil(t, model.Requestors)
	assert.Equal(t, "AND", model.Requestors.LogicalOperator.ValueString())
	require.Len(t, model.Requestors.Conditions, 1)
	assert.Equal(t, "user", model.Requestors.Conditions[0].Type.ValueString())
	assert.Equal(t, "is", model.Requestors.Conditions[0].MatchOperator.ValueString())
	assert.True(t, model.Requestors.Conditions[0].SourceIntegrationName.IsNull())

	var values []string
	diags := model.Requestors.Conditions[0].Values.ElementsAs(ctx, &values, false)
	require.False(t, diags.HasError())
	assert.ElementsMatch(t, []string{"person@example.com"}, values)

	require.Len(t, model.AccessTargets, 1)
	require.NotNil(t, model.AccessTargets[0].Bundle)
	assert.Equal(t, "QA ENV", model.AccessTargets[0].Bundle.Name.ValueString())

	assert.Nil(t, model.RequestFor)
}
