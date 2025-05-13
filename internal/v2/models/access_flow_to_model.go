package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/apono-io/terraform-provider-apono/internal/v2/common"
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

	if val, ok := response.GrantDurationInMin.Get(); ok {
		model.GrantDurationInMin = types.Int32Value(val)
	}

	if val, ok := response.Timeframe.Get(); ok {
		timeframe, err := convertTimeframeToModel(ctx, val)
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

	if val, ok := response.ApproverPolicy.Get(); ok {
		approverPolicy, err := convertApproverPolicyToModel(ctx, val)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver policy: %w", err)
		}
		model.ApproverPolicy = approverPolicy
	}

	requestors, err := convertRequestorsToModel(ctx, response.Requestors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert requestors: %w", err)
	}
	model.Requestors = requestors

	accessTargets, err := convertAccessTargetsToModel(ctx, response.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}
	model.AccessTargets = accessTargets

	return &model, nil
}

func convertTimeframeToModel(ctx context.Context, timeframe client.AccessFlowTimeframePublicV2Model) (*AccessFlowTimeframeModel, error) {
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

func convertApproverPolicyToModel(ctx context.Context, policy client.ApproverPolicyPublicV2Model) (*AccessFlowApproverPolicy, error) {
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

func convertRequestorsToModel(ctx context.Context, requestors client.RequestorsPublicV2Model) (*AccessFlowRequestorsModel, error) {
	model := &AccessFlowRequestorsModel{
		LogicalOperator: types.StringValue(requestors.LogicalOperator),
	}

	var conditions []AccessFlowCondition
	for _, condition := range requestors.Conditions {
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
		Type:   types.StringValue(condition.Type),
		Values: basetypes.NewSetNull(types.StringType),
	}

	if val, ok := condition.SourceIntegrationName.Get(); ok {
		model.SourceIntegrationName = types.StringValue(val)
	}

	if val, ok := condition.MatchOperator.Get(); ok {
		model.MatchOperator = types.StringValue(val)
	} else {
		model.MatchOperator = types.StringValue(common.DefaultMatchOperator)
	}

	if val, ok := condition.Values.Get(); ok {
		valuesSet, diags := types.SetValueFrom(ctx, types.StringType, val)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert condition values: %v", diags)
		}
		model.Values = valuesSet
	}

	return model, nil
}

func convertAccessTargetsToModel(ctx context.Context, accessTargets []client.AccessTargetPublicV2Model) ([]AccessFlowAccessTargetModel, error) {
	var modelTargets []AccessFlowAccessTargetModel

	for _, target := range accessTargets {
		modelTarget := AccessFlowAccessTargetModel{}

		if val, ok := target.Integration.Get(); ok {
			integrationTarget, err := convertIntegrationTargetToModel(ctx, val)
			if err != nil {
				return nil, fmt.Errorf("failed to convert integration target: %w", err)
			}
			modelTarget.Integration = integrationTarget
		}

		if val, ok := target.Bundle.Get(); ok {
			modelTarget.Bundle = &AccessFlowTargetBundleModel{
				Name: types.StringValue(val.BundleName),
			}
		}

		if val, ok := target.AccessScope.Get(); ok {
			modelTarget.AccessScope = &AccessScopeTargetModel{
				Name: types.StringValue(val.AccessScopeName),
			}
		}

		modelTargets = append(modelTargets, modelTarget)
	}

	return modelTargets, nil
}

func convertIntegrationTargetToModel(ctx context.Context, integration client.IntegrationAccessTargetPublicV2Model) (*IntegrationTargetModel, error) {
	model := &IntegrationTargetModel{
		IntegrationName: types.StringValue(integration.IntegrationName),
		ResourceType:    types.StringValue(integration.ResourceType),
	}

	permissionsSet, diags := types.SetValueFrom(ctx, types.StringType, integration.Permissions)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert permissions: %v", diags)
	}
	model.Permissions = permissionsSet

	if val, ok := integration.ResourcesScopes.Get(); ok {
		scopes, err := convertResourcesScopesToModel(ctx, val)
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

		if val, ok := scope.Key.Get(); ok {
			modelScope.Key = types.StringValue(val)
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
		RequesterCannotApproveSelf: types.BoolValue(settings.RequestorCannotApproveHimself),
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
