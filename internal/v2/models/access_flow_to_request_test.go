package models

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/testcommon"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessFlowV2ModelToUpsertRequest(t *testing.T) {
	model := AccessFlowV2Model{
		Name:               types.StringValue("postgresql_prod"),
		Active:             types.BoolValue(true),
		Trigger:            types.StringValue("SELF_SERVE"),
		GrantDurationInMin: types.Int32Null(),
		Timeframe: &AccessFlowTimeframeModel{
			StartTime:  types.StringValue("10:00"),
			EndTime:    types.StringValue("23:59"),
			TimeZone:   types.StringValue("Asia/Jerusalem"),
			DaysOfWeek: testcommon.CreateTestStringSet(t, []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"}),
		},
		Grantees: &AccessFlowGranteesModel{
			LogicalOperator: types.StringValue("OR"),
			Conditions: []AccessFlowCondition{
				{
					SourceIntegrationName: types.StringValue("Okta Directory"),
					Type:                  types.StringValue("user"),
					MatchOperator:         types.StringValue("is"),
					Values:                testcommon.CreateTestStringSet(t, []string{"person@example.com", "person_two@example.com"}),
				},
			},
		},
		AccessTargets: []AccessTargetModel{
			{
				Bundle: &AccessTargetBundleModel{
					Name: types.StringValue("PROD ENV"),
				},
			},
			{
				Integration: &IntegrationTargetModel{
					IntegrationName: types.StringValue("postgresql"),
					ResourceType:    types.StringValue("database"),
					Permissions:     testcommon.CreateTestStringSet(t, []string{"read", "write"}),
					ResourcesScopes: []IntegrationTargetScopeModel{
						{
							ScopeMode: types.StringValue("include_resources"),
							Type:      types.StringValue("NAME"),
							Key:       types.StringNull(),
							Values:    testcommon.CreateTestStringSet(t, []string{"db1", "db2"}),
						},
					},
				},
			},
		},
		ApproverPolicy: &AccessFlowApproverPolicy{
			ApprovalMode: types.StringValue("ANY_OF"),
			ApproverGroups: []AccessFlowApproverGroup{
				{
					LogicalOperator: types.StringValue("OR"),
					Approvers: []AccessFlowCondition{
						{
							SourceIntegrationName: types.StringValue("Okta Directory"),
							Type:                  types.StringValue("user"),
							MatchOperator:         types.StringValue("is"),
							Values:                testcommon.CreateTestStringSet(t, []string{"person@example.com", "person_two@example.com"}),
						},
					},
				},
			},
		},
		Settings: &AccessFlowSettingsModel{
			JustificationRequired:      types.BoolValue(true),
			RequireApproverReason:      types.BoolValue(false),
			RequesterCannotApproveSelf: types.BoolValue(false),
			RequireMFA:                 types.BoolValue(false),
			Labels:                     testcommon.CreateTestStringSet(t, []string{"DB", "PROD", "TERRAFORM"}),
		},
	}

	ctx := t.Context()
	result, err := AccessFlowV2ModelToUpsertRequest(ctx, model)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "postgresql_prod", result.Name)
	assert.True(t, result.Active)
	assert.Equal(t, "SELF_SERVE", result.Trigger)
	assert.False(t, result.GrantDurationInMin.IsSet())

	require.True(t, result.Timeframe.IsSet())
	timeframe, ok := result.Timeframe.Get()
	require.True(t, ok)
	assert.Equal(t, "10:00", timeframe.StartTime)
	assert.Equal(t, "23:59", timeframe.EndTime)
	assert.Equal(t, "Asia/Jerusalem", timeframe.TimeZone)
	assert.ElementsMatch(t,
		[]client.DayOfWeekPublicV2Model{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY"},
		timeframe.DaysOfWeek)

	assert.Equal(t, "OR", result.Grantees.LogicalOperator)
	require.Len(t, result.Grantees.Conditions, 1)
	assert.Equal(t, "user", result.Grantees.Conditions[0].Type)
	sourceIntegRef, ok := result.Grantees.Conditions[0].SourceIntegrationReference.Get()
	require.True(t, ok)
	assert.Equal(t, "Okta Directory", sourceIntegRef)
	matchOp, ok := result.Grantees.Conditions[0].MatchOperator.Get()
	require.True(t, ok)
	assert.Equal(t, "is", matchOp)
	values, ok := result.Grantees.Conditions[0].Values.Get()
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"person@example.com", "person_two@example.com"}, values)

	require.Len(t, result.AccessTargets, 2)

	assert.True(t, result.AccessTargets[0].Bundle.IsSet())
	bundle, ok := result.AccessTargets[0].Bundle.Get()
	require.True(t, ok)
	assert.Equal(t, "PROD ENV", bundle.BundleReference)

	assert.True(t, result.AccessTargets[1].Integration.IsSet())
	integration, ok := result.AccessTargets[1].Integration.Get()
	require.True(t, ok)
	assert.Equal(t, "postgresql", integration.IntegrationReference)
	assert.Equal(t, "database", integration.ResourceType)
	assert.ElementsMatch(t, []string{"read", "write"}, integration.Permissions)
	require.True(t, integration.ResourcesScopes.IsSet())
	resourceScopes, ok := integration.ResourcesScopes.Get()
	require.True(t, ok)
	require.Len(t, resourceScopes, 1)
	assert.Equal(t, "include_resources", resourceScopes[0].ScopeMode)
	assert.Equal(t, "NAME", resourceScopes[0].Type)
	assert.ElementsMatch(t, []string{"db1", "db2"}, resourceScopes[0].Values)

	require.True(t, result.ApproverPolicy.IsSet())
	approverPolicy, ok := result.ApproverPolicy.Get()
	require.True(t, ok)
	assert.Equal(t, "ANY_OF", approverPolicy.ApprovalMode)
	require.Len(t, approverPolicy.ApproverGroups, 1)
	assert.Equal(t, "OR", approverPolicy.ApproverGroups[0].LogicalOperator)
	require.Len(t, approverPolicy.ApproverGroups[0].Approvers, 1)
	assert.Equal(t, "user", approverPolicy.ApproverGroups[0].Approvers[0].Type)
	sourceIntegRef, ok = approverPolicy.ApproverGroups[0].Approvers[0].SourceIntegrationReference.Get()
	require.True(t, ok)
	assert.Equal(t, "Okta Directory", sourceIntegRef)
	matchOp, ok = approverPolicy.ApproverGroups[0].Approvers[0].MatchOperator.Get()
	require.True(t, ok)
	assert.Equal(t, "is", matchOp)
	values, ok = approverPolicy.ApproverGroups[0].Approvers[0].Values.Get()
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"person@example.com", "person_two@example.com"}, values)

	assert.True(t, result.Settings.JustificationRequired)
	assert.False(t, result.Settings.RequireApproverReason)
	assert.False(t, result.Settings.ApproverCannotApproveHimself)
	assert.False(t, result.Settings.RequireMfa)
	assert.ElementsMatch(t, []string{"DB", "PROD", "TERRAFORM"}, result.Settings.Labels)
}

func TestAccessFlowV2ModelToUpsertRequest_NullValues(t *testing.T) {
	model := AccessFlowV2Model{
		Name:    types.StringValue("minimal_flow"),
		Active:  types.BoolValue(false),
		Trigger: types.StringValue("AUTOMATIC"),
		Grantees: &AccessFlowGranteesModel{
			LogicalOperator: types.StringValue("AND"),
			Conditions: []AccessFlowCondition{
				{
					Type:          types.StringValue("user"),
					MatchOperator: types.StringValue("is"),
					Values:        testcommon.CreateTestStringSet(t, []string{"person@example.com"}),
				},
			},
		},
		AccessTargets: []AccessTargetModel{
			{
				Bundle: &AccessTargetBundleModel{
					Name: types.StringValue("QA ENV"),
				},
			},
		},
	}

	ctx := t.Context()
	result, err := AccessFlowV2ModelToUpsertRequest(ctx, model)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "minimal_flow", result.Name)
	assert.False(t, result.Active)
	assert.Equal(t, "AUTOMATIC", result.Trigger)
	assert.False(t, result.GrantDurationInMin.IsSet())
	assert.False(t, result.Timeframe.IsSet())
	assert.False(t, result.ApproverPolicy.IsSet())

	assert.Equal(t, "AND", result.Grantees.LogicalOperator)
	require.Len(t, result.Grantees.Conditions, 1)
	assert.Equal(t, "user", result.Grantees.Conditions[0].Type)
	matchOp, ok := result.Grantees.Conditions[0].MatchOperator.Get()
	require.True(t, ok)
	assert.Equal(t, "is", matchOp)
	values, ok := result.Grantees.Conditions[0].Values.Get()
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"person@example.com"}, values)

	require.Len(t, result.AccessTargets, 1)
	assert.True(t, result.AccessTargets[0].Bundle.IsSet())
	bundle, ok := result.AccessTargets[0].Bundle.Get()
	require.True(t, ok)
	assert.Equal(t, "QA ENV", bundle.BundleReference)
}
