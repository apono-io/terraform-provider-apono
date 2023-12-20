package services

import (
	"context"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertAccessBundleApiToTerraformModel(ctx context.Context, aponoClient *apono.APIClient, accessBundle *apono.AccessBundleV1) (*models.AccessBundleModel, diag.Diagnostics) {
	integrationTargets := accessBundle.GetIntegrationTargets()
	dataIntegrationTargets, diagnostics := ConvertIntegrationTargetsApiToTerraformModel(ctx, aponoClient, integrationTargets)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	accessBundleModel := models.AccessBundleModel{
		ID:                 types.StringValue(accessBundle.GetId()),
		Name:               types.StringValue(accessBundle.GetName()),
		IntegrationTargets: dataIntegrationTargets,
	}

	return &accessBundleModel, nil
}

func ConvertAccessBundleTerraformModelToUpsertApi(ctx context.Context, aponoClient *apono.APIClient, accessBundle *models.AccessBundleModel) (*apono.UpsertAccessBundleV1, diag.Diagnostics) {
	dataIntegrationTargets, diagnostics := ConvertIntegrationTargetsTerraformModelToApi(ctx, aponoClient, accessBundle.IntegrationTargets)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := apono.UpsertAccessBundleV1{
		Name:               accessBundle.Name.ValueString(),
		IntegrationTargets: dataIntegrationTargets,
	}

	return &data, nil
}

func ConvertAccessBundleTerraformModelToUpdateApi(ctx context.Context, aponoClient *apono.APIClient, accessBundle *models.AccessBundleModel) (*apono.UpdateAccessBundleV1, diag.Diagnostics) {
	updateAccessBundleRequest, diagnostics := ConvertAccessBundleTerraformModelToUpsertApi(ctx, aponoClient, accessBundle)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := apono.UpdateAccessBundleV1{
		Name:               *apono.NewNullableString(&updateAccessBundleRequest.Name),
		IntegrationTargets: updateAccessBundleRequest.IntegrationTargets,
	}

	return &data, nil
}
