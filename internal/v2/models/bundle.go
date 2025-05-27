package models

import (
	"context"
	"fmt"

	"github.com/apono-io/terraform-provider-apono/internal/v2/api/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IntegrationTargetModel struct {
	IntegrationName types.String                  `tfsdk:"integration_name"`
	ResourceType    types.String                  `tfsdk:"resource_type"`
	Permissions     types.Set                     `tfsdk:"permissions"`
	ResourcesScopes []IntegrationTargetScopeModel `tfsdk:"resources_scopes"`
}

type IntegrationTargetScopeModel struct {
	ScopeMode types.String `tfsdk:"scope_mode"`
	Type      types.String `tfsdk:"type"`
	Key       types.String `tfsdk:"key"`
	Values    types.List   `tfsdk:"values"`
}

type AccessScopeTargetModel struct {
	Name types.String `tfsdk:"name"`
}

type BundleAccessTargetModel struct {
	Integration *IntegrationTargetModel `tfsdk:"integration"`
	AccessScope *AccessScopeTargetModel `tfsdk:"access_scope"`
}

type BundleV2Model struct {
	ID            types.String              `tfsdk:"id"`
	Name          types.String              `tfsdk:"name"`
	AccessTargets []BundleAccessTargetModel `tfsdk:"access_targets"`
}

type BundlesDataModel struct {
	Name    types.String    `tfsdk:"name"`
	Bundles []BundleV2Model `tfsdk:"bundles"`
}

func BundleResponseToModel(ctx context.Context, response client.BundlePublicV2Model) (*BundleV2Model, error) {
	model := BundleV2Model{
		ID:   types.StringValue(response.ID),
		Name: types.StringValue(response.Name),
	}

	accessTargets, err := convertBundleAccessTargetsToModel(ctx, response.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}
	model.AccessTargets = accessTargets

	return &model, nil
}

func BundlesResponseToModels(ctx context.Context, bundles []client.BundlePublicV2Model) ([]BundleV2Model, error) {
	var bundleModels []BundleV2Model

	for _, bundle := range bundles {
		model, err := BundleResponseToModel(ctx, bundle)
		if err != nil {
			return nil, err
		}

		bundleModels = append(bundleModels, *model)
	}

	return bundleModels, nil
}

func BundleModelToUpsertRequest(ctx context.Context, model BundleV2Model) (*client.UpsertBundlePublicV2Model, error) {
	upsert := client.UpsertBundlePublicV2Model{
		Name: model.Name.ValueString(),
	}

	targets, err := convertBundleAccessTargetsToUpsertRequest(ctx, model.AccessTargets)
	if err != nil {
		return nil, fmt.Errorf("failed to convert access targets: %w", err)
	}

	upsert.AccessTargets = targets

	return &upsert, nil
}

func convertBundleAccessTargetsToModel(ctx context.Context, accessTargets []client.AccessBundleAccessTargetPublicV2Model) ([]BundleAccessTargetModel, error) {
	var modelTargets []BundleAccessTargetModel

	for _, target := range accessTargets {
		modelTarget := BundleAccessTargetModel{}

		if val, ok := target.Integration.Get(); ok {
			integrationTarget, err := convertIntegrationTargetToModel(ctx, val)
			if err != nil {
				return nil, fmt.Errorf("failed to convert integration target: %w", err)
			}
			modelTarget.Integration = integrationTarget
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

func convertBundleAccessTargetsToUpsertRequest(ctx context.Context, models []BundleAccessTargetModel) ([]client.AccessBundleAccessTargetUpsertPublicV2Model, error) {
	var targets []client.AccessBundleAccessTargetUpsertPublicV2Model

	for i, model := range models {
		target := client.AccessBundleAccessTargetUpsertPublicV2Model{}
		setCount := 0

		if model.Integration != nil {
			integration, err := convertIntegrationTargetToUpsertRequest(ctx, *model.Integration)
			if err != nil {
				return nil, fmt.Errorf("failed to convert integration target at index %d: %w", i, err)
			}

			target.Integration.SetTo(*integration)
			setCount++
		}

		if model.AccessScope != nil {
			scope := client.AccessScopeAccessTargetUpsertPublicV2Model{
				AccessScopeReference: model.AccessScope.Name.ValueString(),
			}

			target.AccessScope.SetTo(scope)
			setCount++
		}

		if setCount != 1 {
			return nil, fmt.Errorf("exactly one of 'integration' or 'access_scope' must be configured for each access target (index %d)", i)
		}

		targets = append(targets, target)
	}

	return targets, nil
}
