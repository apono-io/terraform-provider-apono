package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
)

func AccessFlowModelToUpsertRequest(ctx context.Context, model AccessFlowV2Model) (*client.AccessFlowUpsertV2, error) {
	upsert := client.AccessFlowUpsertV2{
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

	requestors, err := convertRequestorsToUpsertRequest(ctx, *model.Requestors)
	if err != nil {
		return nil, fmt.Errorf("failed to convert requestors: %w", err)
	}

	upsert.Requestors = *requestors

	targets, err := convertAccessTargetsToUpsertRequest(ctx, model.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}

	upsert.AccessTargets = targets

	return &upsert, nil
}

func convertTimeframeToUpsertRequest(ctx context.Context, model AccessFlowTimeframeModel) (*client.AccessFlowTimeframeV2, error) {
	timeframe := client.AccessFlowTimeframeV2{
		StartTime: model.StartTime.ValueString(),
		EndTime:   model.EndTime.ValueString(),
		TimeZone:  model.TimeZone.ValueString(),
	}

	var daysOfWeekStrings []string
	if diags := model.DaysOfWeek.ElementsAs(ctx, &daysOfWeekStrings, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert days_of_week: %v", diags)
	}

	daysOfWeek := []client.DayOfWeekV2{}
	for _, dayStr := range daysOfWeekStrings {
		daysOfWeek = append(daysOfWeek, client.DayOfWeekV2(dayStr))
	}

	timeframe.DaysOfWeek = daysOfWeek

	return &timeframe, nil
}

func convertApproverPolicyToUpsertRequest(ctx context.Context, model AccessFlowApproverPolicy) (*client.ApproverPolicyUpsertV2, error) {
	policy := client.ApproverPolicyUpsertV2{
		ApprovalMode: model.ApprovalMode.ValueString(),
	}

	var groups []client.ApproverGroupUpsertV2

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

func convertApproverGroupToUpsertRequest(ctx context.Context, model AccessFlowApproverGroup) (*client.ApproverGroupUpsertV2, error) {
	group := client.ApproverGroupUpsertV2{
		LogicalOperator: model.LogicalOperator.ValueString(),
	}

	var approvers []client.ConditionUpsertV2

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

func convertRequestorsToUpsertRequest(ctx context.Context, model AccessFlowRequestorsModel) (*client.RequestorsUpsertV2, error) {
	requestors := client.RequestorsUpsertV2{
		LogicalOperator: model.LogicalOperator.ValueString(),
	}

	var conditions []client.ConditionUpsertV2

	for i, conditionModel := range model.Conditions {
		condition, err := convertConditionToUpsertRequest(ctx, conditionModel)

		if err != nil {
			return nil, fmt.Errorf("failed to convert grantee condition at index %d: %w", i, err)
		}

		conditions = append(conditions, *condition)

	}

	requestors.Conditions = conditions

	return &requestors, nil
}

func convertConditionToUpsertRequest(ctx context.Context, model AccessFlowCondition) (*client.ConditionUpsertV2, error) {
	condition := client.ConditionUpsertV2{
		Type: model.Type.ValueString(),
	}

	if !model.SourceIntegrationName.IsNull() {
		condition.SourceIntegrationReference.SetTo(model.SourceIntegrationName.ValueString())
	}

	if !model.Values.IsNull() {
		var values []string

		if diags := model.Values.ElementsAs(ctx, &values, false); diags.HasError() {
			return nil, fmt.Errorf("failed to convert values: %v", diags)
		}

		condition.Values.SetTo(values)

		// Set match operator if values are provided.
		condition.MatchOperator.SetTo(model.MatchOperator.ValueString())
	}

	return &condition, nil
}

func convertAccessTargetsToUpsertRequest(ctx context.Context, models []AccessFlowAccessTargetModel) ([]client.AccessTargetUpsertV2, error) {
	var targets []client.AccessTargetUpsertV2

	for i, model := range models {
		target := client.AccessTargetUpsertV2{}
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
			bundle := client.BundleAccessTargetUpsertV2{
				BundleReference: model.Bundle.Name.ValueString(),
			}

			target.Bundle.SetTo(bundle)

			setCount++

		}

		if model.AccessScope != nil {
			scope := client.AccessScopeAccessTargetUpsertV2{
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

func convertIntegrationTargetToUpsertRequest(ctx context.Context, model IntegrationTargetModel) (*client.IntegrationAccessTargetUpsertV2, error) {
	integration := client.IntegrationAccessTargetUpsertV2{
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

func convertResourcesScopesToUpsertRequest(ctx context.Context, scopes []IntegrationTargetScopeModel) ([]client.ResourcesScopeIntegrationAccessTargetV2, error) {
	var result []client.ResourcesScopeIntegrationAccessTargetV2

	for i, scope := range scopes {
		resourceScope := client.ResourcesScopeIntegrationAccessTargetV2{
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

func convertSettingsToUpsertRequest(ctx context.Context, model AccessFlowSettingsModel) (*client.AccessFlowSettingsV2, error) {
	settings := client.AccessFlowSettingsV2{
		JustificationRequired:         model.JustificationRequired.ValueBool(),
		RequireApproverReason:         model.RequireApproverReason.ValueBool(),
		RequestorCannotApproveHimself: model.RequesterCannotApproveSelf.ValueBool(),
		RequireMfa:                    model.RequireMFA.ValueBool(),
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
