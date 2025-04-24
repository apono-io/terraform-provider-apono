package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AccessFlowResponseToModel(ctx context.Context, response client.AccessFlowPublicV2Model) (*AccessFlowV2Model, error) {
	model := AccessFlowV2Model{
		ID:      types.StringValue(response.ID),
		Name:    types.StringValue(response.Name),
		Active:  types.BoolValue(response.Active),
		Trigger: types.StringValue(response.Trigger),
	}

	if response.GrantDurationInMin.IsSet() {
		model.GrantDurationInMin = types.Int32Value(response.GrantDurationInMin.Value)
	}

	if response.Timeframe.IsSet() && !response.Timeframe.IsNull() {
		timeframe, err := convertTimeframeToModel(ctx, response.Timeframe.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert timeframe: %w", err)
		}
		model.Timeframe = timeframe
	}

	settings, err := convertSettingsToModel(ctx, response.Settings)
	if err != nil {
		return nil, fmt.Errorf("failed to convert settings: %w", err)
	}
	model.Settings = settings

	if response.ApproverPolicy.IsSet() && !response.ApproverPolicy.IsNull() {
		approverPolicy, err := convertApproverPolicyToModel(ctx, response.ApproverPolicy.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver policy: %w", err)
		}
		model.ApproverPolicy = approverPolicy
	}

	grantees, err := convertGranteesToModel(ctx, response.Grantees)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grantees: %w", err)
	}
	model.Grantees = grantees

	accessTargets, err := convertAccessTargetsToModel(ctx, response.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}
	model.AccessTargets = accessTargets

	return &model, nil
}

func convertTimeframeToModel(ctx context.Context, timeframe client.AccessFlowPublicV2ModelTimeframe) (*AccessFlowTimeframeModel, error) {
	daysOfWeek := []string{}
	for _, day := range timeframe.DaysOfWeek {
		daysOfWeek = append(daysOfWeek, string(day))
	}

	daysOfWeekSet, diags := types.SetValueFrom(ctx, types.StringType, daysOfWeek)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to create days_of_week set: %v", diags)
	}

	return &AccessFlowTimeframeModel{
		StartTime:  types.StringValue(timeframe.StartTime),
		EndTime:    types.StringValue(timeframe.EndTime),
		DaysOfWeek: daysOfWeekSet,
		TimeZone:   types.StringValue(timeframe.TimeZone),
	}, nil
}

func convertApproverPolicyToModel(ctx context.Context, policy client.AccessFlowPublicV2ModelApproverPolicy) (*AccessFlowApproverPolicy, error) {
	model := &AccessFlowApproverPolicy{
		ApprovalMode: types.StringValue(policy.ApprovalMode),
	}

	var approverGroups []AccessFlowApproverGroup
	for _, group := range policy.ApproverGroups {
		approverGroup, err := convertApproverGroupToModel(ctx, group)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver group: %w", err)
		}
		approverGroups = append(approverGroups, *approverGroup)
	}
	model.ApproverGroups = approverGroups

	return model, nil
}

func convertApproverGroupToModel(ctx context.Context, group client.ApproverGroupPublicV2Model) (*AccessFlowApproverGroup, error) {
	model := &AccessFlowApproverGroup{
		LogicalOperator: types.StringValue(group.LogicalOperator),
	}

	var approvers []AccessFlowCondition
	for _, condition := range group.Approvers {
		approver, err := convertConditionToModel(ctx, condition)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver condition: %w", err)
		}
		approvers = append(approvers, *approver)
	}
	model.Approvers = approvers

	return model, nil
}

func convertGranteesToModel(ctx context.Context, grantees client.GranteesPublicV2Model) (*AccessFlowGranteesModel, error) {
	model := &AccessFlowGranteesModel{
		LogicalOperator: types.StringValue(grantees.LogicalOperator),
	}

	var conditions []AccessFlowCondition
	for _, condition := range grantees.Conditions {
		conditionModel, err := convertConditionToModel(ctx, condition)
		if err != nil {
			return nil, fmt.Errorf("failed to convert grantee condition: %w", err)
		}
		conditions = append(conditions, *conditionModel)
	}
	model.Conditions = conditions

	return model, nil
}

func convertConditionToModel(ctx context.Context, condition client.ConditionPublicV2Model) (*AccessFlowCondition, error) {
	model := &AccessFlowCondition{
		Type: types.StringValue(condition.Type),
	}

	if condition.SourceIntegrationName.IsSet() && !condition.SourceIntegrationName.IsNull() {
		model.SourceIntegrationName = types.StringValue(condition.SourceIntegrationName.Value)
	}

	if condition.MatchOperator.IsSet() && !condition.MatchOperator.IsNull() {
		model.MatchOperator = types.StringValue(condition.MatchOperator.Value)
	}

	if condition.Values.IsSet() && !condition.Values.IsNull() {
		valuesSet, diags := types.SetValueFrom(ctx, types.StringType, condition.Values.Value)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert condition values: %v", diags)
		}
		model.Values = valuesSet
	}

	return model, nil
}

func convertAccessTargetsToModel(ctx context.Context, accessTargets []client.AccessTargetPublicV2Model) ([]AccessTargetModel, error) {
	var modelTargets []AccessTargetModel

	for _, target := range accessTargets {
		modelTarget := AccessTargetModel{}

		if target.Integration.IsSet() && !target.Integration.IsNull() {
			integrationTarget, err := convertIntegrationTargetToModel(ctx, target.Integration.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert integration target: %w", err)
			}
			modelTarget.Integration = integrationTarget
		}

		if target.Bundle.IsSet() && !target.Bundle.IsNull() {
			modelTarget.Bundle = &AccessTargetBundleModel{
				Name: types.StringValue(target.Bundle.Value.BundleName),
			}
		}

		if target.AccessScope.IsSet() && !target.AccessScope.IsNull() {
			modelTarget.AccessScope = &AccessScopeTargetModel{
				Name: types.StringValue(target.AccessScope.Value.AccessScopeName),
			}
		}

		modelTargets = append(modelTargets, modelTarget)
	}

	return modelTargets, nil
}

func convertIntegrationTargetToModel(ctx context.Context, integration client.AccessTargetPublicV2ModelIntegration) (*IntegrationTargetModel, error) {
	model := &IntegrationTargetModel{
		IntegrationName: types.StringValue(integration.IntegrationName),
		ResourceType:    types.StringValue(integration.ResourceType),
	}

	permissionsSet, diags := types.SetValueFrom(ctx, types.StringType, integration.Permissions)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert permissions: %v", diags)
	}
	model.Permissions = permissionsSet

	if integration.ResourcesScopes.IsSet() && !integration.ResourcesScopes.IsNull() {
		scopes, err := convertResourcesScopesToModel(ctx, integration.ResourcesScopes.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource scopes: %w", err)
		}
		model.ResourcesScopes = scopes
	}

	return model, nil
}

func convertResourcesScopesToModel(ctx context.Context, scopes []client.ResourcesScopeIntegrationAccessTargetPublicV2Model) ([]IntegrationTargetScopeModel, error) {
	var modelScopes []IntegrationTargetScopeModel

	for _, scope := range scopes {
		modelScope := IntegrationTargetScopeModel{
			ScopeMode: types.StringValue(scope.ScopeMode),
			Type:      types.StringValue(scope.Type),
		}

		if scope.Key.IsSet() {
			modelScope.Key = types.StringValue(scope.Key.Value)
		}

		valuesSet, diags := types.SetValueFrom(ctx, types.StringType, scope.Values)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert resource scope values: %v", diags)
		}
		modelScope.Values = valuesSet

		modelScopes = append(modelScopes, modelScope)
	}

	return modelScopes, nil
}

func convertSettingsToModel(ctx context.Context, settings client.AccessFlowSettingsPublicV2Model) (*AccessFlowSettingsModel, error) {
	model := &AccessFlowSettingsModel{
		JustificationRequired:      types.BoolValue(settings.JustificationRequired),
		RequireApproverReason:      types.BoolValue(settings.RequireApproverReason),
		RequesterCannotApproveSelf: types.BoolValue(settings.ApproverCannotApproveHimself),
		RequireMFA:                 types.BoolValue(settings.RequireMfa),
	}

	if len(settings.Labels) > 0 {
		labelsSet, diags := types.SetValueFrom(ctx, types.StringType, settings.Labels)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert labels: %v", diags)
		}
		model.Labels = labelsSet
	} else {
		model.Labels = basetypes.NewSetNull(types.StringType)
	}

	return model, nil
}
