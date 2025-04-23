package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func AccessFlowV2ModelToUpsertRequest(ctx context.Context, model AccessFlowV2Model) (*client.AccessFlowUpsertPublicV2Model, error) {
	upsert := client.AccessFlowUpsertPublicV2Model{
		Name:    model.Name.ValueString(),
		Active:  model.Active.ValueBool(),
		Trigger: model.Trigger.ValueString(),
	}

	if !model.GrantDurationInMin.IsNull() {
		upsert.GrantDurationInMin.SetTo(model.GrantDurationInMin.ValueInt32())
	}

	if model.Timeframe != nil {
		timeframe, err := convertTimeframeToUpsertRequest(ctx, *model.Timeframe)
		if err != nil {
			return nil, fmt.Errorf("failed to convert timeframe: %w", err)
		}

		upsert.Timeframe.SetTo(*timeframe)
	}

	if model.Settings != nil {
		settings, err := convertSettingsToUpsertRequest(ctx, *model.Settings)
		if err != nil {
			return nil, fmt.Errorf("failed to convert settings: %w", err)
		}

		upsert.Settings = *settings

	}

	if model.ApproverPolicy != nil {
		approverPolicy, err := convertApproverPolicyToUpsertRequest(ctx, *model.ApproverPolicy)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver policy: %w", err)
		}

		upsert.ApproverPolicy.SetTo(*approverPolicy)
	}

	grantees, err := convertGranteesToUpsertRequest(ctx, *model.Grantees)
	if err != nil {
		return nil, fmt.Errorf("failed to convert grantees: %w", err)
	}

	upsert.Grantees = *grantees

	targets, err := convertAccessTargetsToUpsertRequest(ctx, model.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}

	upsert.AccessTargets = targets

	return &upsert, nil
}

func convertTimeframeToUpsertRequest(ctx context.Context, model AccessFlowTimeframeModel) (*client.AccessFlowUpsertPublicV2ModelTimeframe, error) {
	timeframe := client.AccessFlowUpsertPublicV2ModelTimeframe{
		StartTime: model.StartTime.ValueString(),
		EndTime:   model.EndTime.ValueString(),
		TimeZone:  model.TimeZone.ValueString(),
	}

	var daysOfWeekStrings []string
	if diags := model.DaysOfWeek.ElementsAs(ctx, &daysOfWeekStrings, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert days_of_week: %v", diags)
	}

	daysOfWeek := []client.DayOfWeekPublicV2Model{}
	for _, dayStr := range daysOfWeekStrings {
		daysOfWeek = append(daysOfWeek, client.DayOfWeekPublicV2Model(dayStr))
	}

	timeframe.DaysOfWeek = daysOfWeek

	return &timeframe, nil
}

func convertApproverPolicyToUpsertRequest(ctx context.Context, model AccessFlowApproverPolicy) (*client.AccessFlowUpsertPublicV2ModelApproverPolicy, error) {
	policy := client.AccessFlowUpsertPublicV2ModelApproverPolicy{
		ApprovalMode: model.ApprovalMode.ValueString(),
	}

	var groups []client.ApproverGroupUpsertPublicV2Model

	for i, groupModel := range model.ApproverGroups {
		group, err := convertApproverGroupToUpsertRequest(ctx, groupModel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver group at index %d: %w", i, err)
		}

		groups = append(groups, *group)
	}

	policy.ApproverGroups = groups

	return &policy, nil
}

func convertApproverGroupToUpsertRequest(ctx context.Context, model AccessFlowApproverGroup) (*client.ApproverGroupUpsertPublicV2Model, error) {
	group := client.ApproverGroupUpsertPublicV2Model{
		LogicalOperator: model.LogicalOperator.ValueString(),
	}

	var approvers []client.ConditionUpsertPublicV2Model

	for i, approverModel := range model.Approvers {
		condition, err := convertConditionToUpsertRequest(ctx, approverModel)
		if err != nil {
			return nil, fmt.Errorf("failed to convert approver condition at index %d: %w", i, err)
		}

		approvers = append(approvers, *condition)
	}

	group.Approvers = approvers

	return &group, nil
}

func convertGranteesToUpsertRequest(ctx context.Context, model AccessFlowGranteesModel) (*client.GranteesUpsertPublicV2Model, error) {
	grantees := client.GranteesUpsertPublicV2Model{
		LogicalOperator: model.LogicalOperator.ValueString(),
	}

	var conditions []client.ConditionUpsertPublicV2Model

	for i, conditionModel := range model.Conditions {
		condition, err := convertConditionToUpsertRequest(ctx, conditionModel)

		if err != nil {
			return nil, fmt.Errorf("failed to convert grantee condition at index %d: %w", i, err)
		}

		conditions = append(conditions, *condition)

	}

	grantees.Conditions = conditions

	return &grantees, nil
}

func convertConditionToUpsertRequest(ctx context.Context, model AccessFlowCondition) (*client.ConditionUpsertPublicV2Model, error) {
	condition := client.ConditionUpsertPublicV2Model{
		Type: model.Type.ValueString(),
	}

	if !model.SourceIntegrationName.IsNull() {
		condition.SourceIntegrationReference.SetTo(model.SourceIntegrationName.ValueString())
	}

	condition.MatchOperator.SetTo(model.MatchOperator.ValueString())

	var values []string

	if diags := model.Values.ElementsAs(ctx, &values, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert values: %v", diags)
	}

	condition.Values.SetTo(values)

	return &condition, nil
}

func convertAccessTargetsToUpsertRequest(ctx context.Context, models []AccessTargetModel) ([]client.AccessTargetUpsertPublicV2Model, error) {
	var targets []client.AccessTargetUpsertPublicV2Model

	for i, model := range models {
		target := client.AccessTargetUpsertPublicV2Model{}
		setCount := 0

		if model.Integration != nil {
			integration, err := convertIntegrationTargetToUpsertRequest(ctx, *model.Integration)
			if err != nil {
				return nil, fmt.Errorf("failed to convert integration target at index %d: %w", i, err)
			}

			target.Integration.SetTo(*integration)

			setCount++
		}

		if model.Bundle != nil {
			bundle := client.AccessTargetUpsertPublicV2ModelBundle{
				BundleReference: model.Bundle.Name.ValueString(),
			}

			target.Bundle.SetTo(bundle)

			setCount++

		}

		if model.AccessScope != nil {
			scope := client.AccessTargetUpsertPublicV2ModelAccessScope{
				AccessScopeReference: model.AccessScope.Name.ValueString(),
			}

			target.AccessScope.SetTo(scope)

			setCount++
		}

		if setCount != 1 {
			return nil, fmt.Errorf("exactly one of 'integration', 'bundle', or 'access_scope' must be configured for each access target (index %d)", i)
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func convertIntegrationTargetToUpsertRequest(ctx context.Context, model IntegrationTargetModel) (*client.AccessTargetUpsertPublicV2ModelIntegration, error) {
	integration := client.AccessTargetUpsertPublicV2ModelIntegration{
		IntegrationReference: model.IntegrationName.ValueString(),
		ResourceType:         model.ResourceType.ValueString(),
	}

	var permissions []string
	if diags := model.Permissions.ElementsAs(ctx, &permissions, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert permissions: %v", diags)
	}
	integration.Permissions = permissions

	if len(model.ResourcesScopes) > 0 {
		scopes, err := convertResourcesScopesToUpsertRequest(ctx, model.ResourcesScopes)
		if err != nil {
			return nil, fmt.Errorf("failed to convert resource scopes: %w", err)
		}
		integration.ResourcesScopes.SetTo(scopes)
	}

	return &integration, nil
}

func convertResourcesScopesToUpsertRequest(ctx context.Context, scopes []IntegrationTargetScopeModel) ([]client.ResourcesScopeIntegrationAccessTargetPublicV2Model, error) {
	var result []client.ResourcesScopeIntegrationAccessTargetPublicV2Model

	for i, scope := range scopes {
		resourceScope := client.ResourcesScopeIntegrationAccessTargetPublicV2Model{
			ScopeMode: scope.ScopeMode.ValueString(),
			Type:      scope.Type.ValueString(),
		}

		if !scope.Key.IsNull() {
			resourceScope.Key.SetTo(scope.Key.ValueString())
		}

		var values []string

		if diags := scope.Values.ElementsAs(ctx, &values, false); diags.HasError() {
			return nil, fmt.Errorf("failed to convert resource scope values at index %d: %v", i, diags)
		}

		resourceScope.Values = values

		result = append(result, resourceScope)
	}

	return result, nil
}

func convertSettingsToUpsertRequest(ctx context.Context, model AccessFlowSettingsModel) (*client.AccessFlowSettingsPublicV2Model, error) {
	settings := client.AccessFlowSettingsPublicV2Model{
		JustificationRequired:        model.JustificationRequired.ValueBool(),
		RequireApproverReason:        model.RequireApproverReason.ValueBool(),
		ApproverCannotApproveHimself: model.RequesterCannotApproveSelf.ValueBool(),
		RequireMfa:                   model.RequireMFA.ValueBool(),
	}

	if !model.Labels.IsNull() {
		var labels []string

		if diags := model.Labels.ElementsAs(ctx, &labels, false); diags.HasError() {
			return nil, fmt.Errorf("failed to convert labels: %v", diags)
		}

		settings.Labels = labels
	} else {
		settings.Labels = []string{}
	}

	return &settings, nil
}
