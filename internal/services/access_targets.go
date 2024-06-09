package services

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

func convertIntegrationTargetsApiToTerraformModel(ctx context.Context, aponoClient *apono.APIClient, integrationTargets []apono.AccessTargetIntegrationV1) ([]models.IntegrationTarget, diag.Diagnostics) {
	availableIntegrations, _, err := aponoClient.IntegrationsApi.ListIntegrationsV2(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "integrations", "")
	}

	var dataIntegrationTargets []models.IntegrationTarget
	for _, integrationTarget := range integrationTargets {
		integration, diagnostics := convertIntegrationTargetApiToTerraformModel(ctx, &integrationTarget, availableIntegrations.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataIntegrationTargets = append(dataIntegrationTargets, *integration)
	}

	return dataIntegrationTargets, nil
}

func convertIntegrationTargetsTerraformModelToApi(ctx context.Context, aponoClient *apono.APIClient, integrationTargets []models.IntegrationTarget) ([]apono.AccessTargetIntegrationV1, diag.Diagnostics) {
	availableIntegrations, _, err := aponoClient.IntegrationsApi.ListIntegrationsV2(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "integrations", "")
	}

	var resultIntegrationTargets []apono.AccessTargetIntegrationV1
	for _, integrationTarget := range integrationTargets {
		integration, diagnostics := convertIntegrationTargetTerraformModelToApi(integrationTarget, availableIntegrations.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		resultIntegrationTargets = append(resultIntegrationTargets, *integration)
	}

	return resultIntegrationTargets, nil
}

func convertBundleTargetsApiToTerraformModel(ctx context.Context, aponoClient *apono.APIClient, bundleTargets []aponoapi.AccessTargetBundleTerraformV1) ([]models.BundleTarget, diag.Diagnostics) {
	availableBundles, _, err := aponoClient.AccessBundlesApi.ListAccessBundles(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "access bundles", "")
	}

	var dataBundleTargets []models.BundleTarget
	for _, bundleTarget := range bundleTargets {
		bundle, diagnostics := convertBundleTargetApiToTerraformModel(&bundleTarget, availableBundles.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataBundleTargets = append(dataBundleTargets, *bundle)
	}

	return dataBundleTargets, nil
}

func convertBundleTargetsTerraformModelToApi(ctx context.Context, aponoClient *apono.APIClient, bundleTargets []models.BundleTarget) ([]aponoapi.AccessTargetBundleTerraformV1, diag.Diagnostics) {
	availableBundles, _, err := aponoClient.AccessBundlesApi.ListAccessBundles(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "access bundles", "")
	}

	var resultBundleTargets []aponoapi.AccessTargetBundleTerraformV1
	for _, bundleTarget := range bundleTargets {
		bundle, diagnostics := convertBundleTargetTerraformModelToApi(bundleTarget, availableBundles.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		resultBundleTargets = append(resultBundleTargets, *bundle)
	}

	return resultBundleTargets, nil
}

func convertIntegrationTargetApiToTerraformModel(ctx context.Context, integrationTarget *apono.AccessTargetIntegrationV1, availableIntegrations []apono.Integration) (*models.IntegrationTarget, diag.Diagnostics) {
	var result *models.IntegrationTarget
	for _, integration := range availableIntegrations {
		if integration.Id == integrationTarget.GetIntegrationId() {
			resourceIncludeFilters := convertTagV1ListToResourceFilter(integrationTarget.GetResourceTagIncludes())
			resourceExcludeFilters := convertTagV1ListToResourceFilter(integrationTarget.GetResourceTagExcludes())

			permissions, diagnostics := types.SetValueFrom(ctx, types.StringType, integrationTarget.GetPermissions())
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}

			result = &models.IntegrationTarget{
				Name:                   types.StringValue(integration.GetName()),
				ResourceType:           types.StringValue(integrationTarget.GetResourceType()),
				ResourceIncludeFilters: resourceIncludeFilters,
				ResourceExcludeFilters: resourceExcludeFilters,
				Permissions:            permissions,
			}
		}
	}

	if result == nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration: %s", integrationTarget.GetIntegrationId()))
		return nil, diagnostics
	}

	return result, nil
}

func convertIntegrationTargetTerraformModelToApi(integrationTarget models.IntegrationTarget, availableIntegrations []apono.Integration) (*apono.AccessTargetIntegrationV1, diag.Diagnostics) {
	var result *apono.AccessTargetIntegrationV1
	for _, integration := range availableIntegrations {
		if integration.Name == integrationTarget.Name.ValueString() && slices.Contains(integration.ConnectedResourceTypes, integrationTarget.ResourceType.ValueString()) {
			resourceTagInclude, diagnostics := convertResourceFilterListToTagV1Api(integrationTarget.ResourceIncludeFilters)
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}
			resourceTagExclude, diagnostics := convertResourceFilterListToTagV1Api(integrationTarget.ResourceExcludeFilters)
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}

			var permissions []string
			for _, permission := range integrationTarget.Permissions.Elements() {
				permissions = append(permissions, utils.AttrValueToString(permission))
			}

			result = &apono.AccessTargetIntegrationV1{
				IntegrationId:       integration.Id,
				ResourceType:        integrationTarget.ResourceType.ValueString(),
				ResourceTagIncludes: resourceTagInclude,
				ResourceTagExcludes: resourceTagExclude,
				Permissions:         permissions,
			}
		}
	}

	if result == nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration: (%s) with resource type (%s)", integrationTarget.Name.ValueString(), integrationTarget.ResourceType.ValueString()))
		return nil, diagnostics
	}

	return result, nil
}

func convertBundleTargetApiToTerraformModel(bundleTarget *aponoapi.AccessTargetBundleTerraformV1, availableBundles []apono.AccessBundleV1) (*models.BundleTarget, diag.Diagnostics) {
	var result *models.BundleTarget
	for _, bundle := range availableBundles {
		if bundle.Id == bundleTarget.GetBundleId() {
			result = &models.BundleTarget{
				Name: types.StringValue(bundle.GetName()),
			}
		}
	}

	if result == nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get bundle: %s", bundleTarget.GetBundleId()))
		return nil, diagnostics
	}

	return result, nil
}

func convertBundleTargetTerraformModelToApi(bundleTarget models.BundleTarget, availableBundles []apono.AccessBundleV1) (*aponoapi.AccessTargetBundleTerraformV1, diag.Diagnostics) {
	var result *aponoapi.AccessTargetBundleTerraformV1
	for _, bundle := range availableBundles {
		if bundle.Name == bundleTarget.Name.ValueString() {
			result = &aponoapi.AccessTargetBundleTerraformV1{
				BundleId: bundle.Id,
			}
		}
	}

	if result == nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get bundle: (%s)", bundleTarget.Name.ValueString()))
		return nil, diagnostics
	}

	return result, nil
}

func convertIntegrationTargetsOldApiToNewApiModel(integrationTarget []apono.AccessTargetIntegrationV1) []aponoapi.AccessTargetIntegrationTerraformV1 {
	var result []aponoapi.AccessTargetIntegrationTerraformV1
	for _, target := range integrationTarget {
		result = append(result, aponoapi.AccessTargetIntegrationTerraformV1{
			IntegrationId:       target.GetIntegrationId(),
			ResourceType:        target.GetResourceType(),
			ResourceTagIncludes: convertResourceTagsOldApiToNewApiModel(target.GetResourceTagIncludes()),
			ResourceTagExcludes: convertResourceTagsOldApiToNewApiModel(target.GetResourceTagExcludes()),
			Permissions:         target.GetPermissions(),
		})
	}

	return result
}

func convertIntegrationTargetsNewApiToOldApiModel(integrationTarget []aponoapi.AccessTargetIntegrationTerraformV1) []apono.AccessTargetIntegrationV1 {
	var result []apono.AccessTargetIntegrationV1
	for _, target := range integrationTarget {
		result = append(result, apono.AccessTargetIntegrationV1{
			IntegrationId:       target.IntegrationId,
			ResourceType:        target.ResourceType,
			ResourceTagIncludes: convertResourceTagsNewApiToOldApiModel(target.ResourceTagIncludes),
			ResourceTagExcludes: convertResourceTagsNewApiToOldApiModel(target.ResourceTagExcludes),
			Permissions:         target.Permissions,
		})
	}

	return result
}

func convertResourceTagsOldApiToNewApiModel(resourceTags []apono.TagV1) []aponoapi.TagTerraformV1 {
	result := []aponoapi.TagTerraformV1{}
	for _, resourceTag := range resourceTags {
		result = append(result, aponoapi.TagTerraformV1{
			Name:  resourceTag.GetName(),
			Value: resourceTag.GetValue(),
		})
	}

	return result
}

func convertResourceTagsNewApiToOldApiModel(resourceTags []aponoapi.TagTerraformV1) []apono.TagV1 {
	result := []apono.TagV1{}
	for _, resourceTag := range resourceTags {
		result = append(result, apono.TagV1{
			Name:  resourceTag.Name,
			Value: resourceTag.Value,
		})
	}

	return result
}
