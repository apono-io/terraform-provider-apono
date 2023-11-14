package services

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/apono-io/terraform-provider-apono/internal/models"
	"github.com/apono-io/terraform-provider-apono/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
	"math/big"
	"strings"
)

func ConvertToAccessFlowModel(ctx context.Context, aponoClient *apono.APIClient, accessFlow *apono.AccessFlowV1) (*models.AccessFlowModel, diag.Diagnostics) {
	data := models.AccessFlowModel{}
	data.ID = types.StringValue(accessFlow.GetId())
	data.Name = types.StringValue(accessFlow.GetName())
	data.Active = types.BoolValue(accessFlow.GetActive())
	data.RevokeAfterInSec = types.NumberValue(big.NewFloat(float64(accessFlow.GetRevokeAfterInSec())))

	trigger := accessFlow.GetTrigger()
	dataTrigger := models.Trigger{
		Type: types.StringValue(trigger.GetType()),
	}

	timeframe, ok := trigger.GetTimeframeOk()
	if ok && timeframe != nil {
		var timeframeDaysAsStrings []types.String
		for _, day := range timeframe.GetDaysInWeek() {
			timeframeDaysAsStrings = append(timeframeDaysAsStrings, types.StringValue(string(day)))
		}
		daysInWeek, diagnostics := types.SetValueFrom(ctx, types.StringType, timeframeDaysAsStrings)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataTrigger.Timeframe = &models.Timeframe{
			StartOfDayTimeInSeconds: types.NumberValue(big.NewFloat(float64(timeframe.GetStartOfDayTimeInSeconds()))),
			EndOfDayTimeInSeconds:   types.NumberValue(big.NewFloat(float64(timeframe.GetEndOfDayTimeInSeconds()))),
			DaysInWeek:              daysInWeek,
			TimeZone:                types.StringValue(timeframe.GetTimeZone()),
		}

	}
	data.Trigger = &dataTrigger

	availableIdentities, _, err := aponoClient.IdentitiesApi.ListIdentities(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "identities", "")
	}

	var dataGrantees []models.Identity
	for _, grantee := range accessFlow.GetGrantees() {
		identity, diagnostics := convertToIdentityModel(ctx, grantee.Id, strings.ToLower(grantee.Type), aponoClient, availableIdentities.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataGrantees = append(dataGrantees, *identity)
	}
	setGrantees, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: models.IdentityObject}, getUniqueListOfIdentities(dataGrantees))
	if len(diags) > 0 {
		return nil, diags
	}
	data.Grantees = setGrantees

	var dataApprovers []models.Identity
	for _, approver := range accessFlow.GetApprovers() {
		identity, diagnostics := convertToIdentityModel(ctx, approver.Id, strings.ToLower(approver.Type), aponoClient, availableIdentities.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataApprovers = append(dataApprovers, *identity)
	}
	setApprovers, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: models.IdentityObject}, getUniqueListOfIdentities(dataApprovers))
	if len(diags) > 0 {
		return nil, diags
	}
	data.Approvers = setApprovers

	availableIntegrations, _, err := aponoClient.IntegrationsApi.ListIntegrationsV2(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "integrations", "")
	}

	integrationTargets := accessFlow.GetIntegrationTargets()
	var dataIntegrationTargets []models.IntegrationTarget
	for _, integrationTarget := range integrationTargets {
		integration, diagnostics := convertToIntegrationTargetModel(ctx, &integrationTarget, availableIntegrations.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataIntegrationTargets = append(dataIntegrationTargets, *integration)
	}
	data.IntegrationTargets = dataIntegrationTargets

	settings, ok := accessFlow.GetSettingsOk()
	if ok && settings != nil {
		dataSettings := models.Settings{
			RequireJustificationOnRequestAgain: types.BoolValue(settings.GetRequireJustificationOnRequestAgain()),
			RequireAllApprovers:                types.BoolValue(settings.GetRequireAllApprovers()),
			ApproverCannotApproveHimself:       types.BoolValue(settings.GetApproverCannotApproveHimself()),
		}
		data.Settings = &dataSettings
	}

	return &data, nil
}

func ConvertToAccessFlowUpsertApiModel(ctx context.Context, aponoClient *apono.APIClient, accessFlow *models.AccessFlowModel) (*apono.UpsertAccessFlowV1, diag.Diagnostics) {
	var data apono.UpsertAccessFlowV1

	data.Name = accessFlow.Name.ValueString()
	data.Active = accessFlow.Active.ValueBool()
	revokeAfterInSec, _ := accessFlow.RevokeAfterInSec.ValueBigFloat().Int64()
	data.RevokeAfterInSec = int32(revokeAfterInSec)

	dataTrigger, diagnostics := convertTriggerToApiModel(*accessFlow.Trigger)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}
	data.Trigger = *dataTrigger

	availableIdentities, _, err := aponoClient.IdentitiesApi.ListIdentities(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "identities", "")
	}

	existingGrantees := make([]models.Identity, 0, len(accessFlow.Grantees.Elements()))
	diagnostics = accessFlow.Grantees.ElementsAs(ctx, &existingGrantees, false)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var dataGrantees []apono.GranteeV1
	for _, grantee := range existingGrantees {
		granteeIds, diagnostics := getIdentitiesIdsByNameAndType(ctx, grantee.Name.ValueString(), grantee.Type.ValueString(), availableIdentities.Data, aponoClient)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		for _, granteeId := range granteeIds {
			dataGrantees = append(dataGrantees, apono.GranteeV1{
				Id:   granteeId,
				Type: grantee.Type.ValueString(),
			})
		}
	}
	data.Grantees = dataGrantees

	existingApprovers := make([]models.Identity, 0, len(accessFlow.Approvers.Elements()))
	diagnostics = accessFlow.Approvers.ElementsAs(ctx, &existingApprovers, false)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var dataApprovers []apono.ApproverV1
	for _, approver := range existingApprovers {
		approverIds, diagnostics := getIdentitiesIdsByNameAndType(ctx, approver.Name.ValueString(), approver.Type.ValueString(), availableIdentities.Data, aponoClient)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		for _, approverId := range approverIds {
			dataApprovers = append(dataApprovers, apono.ApproverV1{
				Id:   approverId,
				Type: approver.Type.ValueString(),
			})
		}
	}
	data.Approvers = dataApprovers

	availableIntegrations, _, err := aponoClient.IntegrationsApi.ListIntegrationsV2(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "integrations", "")
	}

	var dataIntegrationTargets []apono.AccessTargetIntegrationV1
	for _, integrationTarget := range accessFlow.IntegrationTargets {
		dataIntegrationTarget, diagnostics := convertIntegrationTargetToApiModel(integrationTarget, availableIntegrations.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataIntegrationTargets = append(dataIntegrationTargets, *dataIntegrationTarget)
	}
	data.IntegrationTargets = dataIntegrationTargets

	if accessFlow.Settings != nil {
		settings := apono.AccessFlowV1Settings{
			RequireJustificationOnRequestAgain: *apono.NewNullableBool(accessFlow.Settings.RequireJustificationOnRequestAgain.ValueBoolPointer()),
			RequireAllApprovers:                *apono.NewNullableBool(accessFlow.Settings.RequireAllApprovers.ValueBoolPointer()),
			ApproverCannotApproveHimself:       *apono.NewNullableBool(accessFlow.Settings.ApproverCannotApproveHimself.ValueBoolPointer()),
		}

		data.Settings = *apono.NewNullableAccessFlowV1Settings(&settings)
	}

	return &data, nil
}

func ConvertToAccessFlowUpdateApiModel(ctx context.Context, aponoClient *apono.APIClient, accessFlow *models.AccessFlowModel) (*apono.UpdateAccessFlowV1, diag.Diagnostics) {
	updateAccessFlowRequest, diagnostics := ConvertToAccessFlowUpsertApiModel(ctx, aponoClient, accessFlow)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	trigger := apono.UpdateAccessFlowV1Trigger{
		Type:      updateAccessFlowRequest.Trigger.Type,
		Timeframe: updateAccessFlowRequest.Trigger.Timeframe,
	}

	data := apono.UpdateAccessFlowV1{
		Name:               *apono.NewNullableString(&updateAccessFlowRequest.Name),
		Active:             *apono.NewNullableBool(&updateAccessFlowRequest.Active),
		RevokeAfterInSec:   *apono.NewNullableInt32(&updateAccessFlowRequest.RevokeAfterInSec),
		Trigger:            *apono.NewNullableUpdateAccessFlowV1Trigger(&trigger),
		Grantees:           updateAccessFlowRequest.Grantees,
		Approvers:          updateAccessFlowRequest.Approvers,
		IntegrationTargets: updateAccessFlowRequest.IntegrationTargets,
		Settings:           updateAccessFlowRequest.Settings,
	}

	return &data, nil

}

func convertToIntegrationTargetModel(ctx context.Context, integrationTarget *apono.AccessTargetIntegrationV1, availableIntegrations []apono.Integration) (*models.IntegrationTarget, diag.Diagnostics) {
	for _, integration := range availableIntegrations {
		if integration.Id == integrationTarget.GetIntegrationId() {
			resourceIncludeFilters := convertTagsToFiltersModel(integrationTarget.GetResourceTagIncludes())
			resourceExcludeFilters := convertTagsToFiltersModel(integrationTarget.GetResourceTagExcludes())

			permissions, diagnostics := types.SetValueFrom(ctx, types.StringType, integrationTarget.GetPermissions())
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}

			return &models.IntegrationTarget{
				Name:                   types.StringValue(integration.GetName()),
				ResourceType:           types.StringValue(integrationTarget.GetResourceType()),
				ResourceIncludeFilters: resourceIncludeFilters,
				ResourceExcludeFilters: resourceExcludeFilters,
				Permissions:            permissions,
			}, nil
		}
	}

	diagnostics := diag.Diagnostics{}
	diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration: %s", integrationTarget.GetIntegrationId()))
	return nil, diagnostics

}

func convertToIdentityModel(ctx context.Context, identityId string, identityType string, aponoClient *apono.APIClient, availableIdentities []apono.IdentityModel2) (*models.Identity, diag.Diagnostics) {
	switch identityType {
	case "user":
		user, _, err := aponoClient.UsersApi.GetUser(ctx, identityId).Execute()
		if err != nil {
			return nil, utils.GetDiagnosticsForApiError(err, "list", "user", identityId)
		}
		return &models.Identity{
			Name: types.StringValue(user.GetEmail()),
			Type: types.StringValue("user"),
		}, nil

	case "group", "context_attribute":
		for _, identity := range availableIdentities {
			if identity.Id == identityId {
				return &models.Identity{
					Name: types.StringValue(identity.GetName()),
					Type: types.StringValue(identityType),
				}, nil
			}
		}

		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get identity %s: %s", identityType, identityId))
		return nil, diagnostics

	default:
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unsupported indentity type: %s, please update the provider to support this type or removed it using UI/API", identityType),
		)
		return nil, diagnostics
	}
}

func convertIntegrationTargetToApiModel(integrationTarget models.IntegrationTarget, availableIntegrations []apono.Integration) (*apono.AccessTargetIntegrationV1, diag.Diagnostics) {
	for _, integration := range availableIntegrations {
		if integration.Name == integrationTarget.Name.ValueString() && slices.Contains(integration.ConnectedResourceTypes, integrationTarget.ResourceType.ValueString()) {
			resourceTagInclude, diagnostics := convertFiltersToListTagsApiModel(integrationTarget.ResourceIncludeFilters)
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}
			resourceTagExclude, diagnostics := convertFiltersToListTagsApiModel(integrationTarget.ResourceExcludeFilters)
			if len(diagnostics) > 0 {
				return nil, diagnostics
			}

			var permissions []string
			for _, permission := range integrationTarget.Permissions.Elements() {
				permissions = append(permissions, utils.AttrValueToString(permission))
			}

			return &apono.AccessTargetIntegrationV1{
				IntegrationId:       integration.Id,
				ResourceType:        integrationTarget.ResourceType.ValueString(),
				ResourceTagIncludes: resourceTagInclude,
				ResourceTagExcludes: resourceTagExclude,
				Permissions:         permissions,
			}, nil
		}
	}

	diagnostics := diag.Diagnostics{}
	diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get integration: (%s) with resource type (%s)", integrationTarget.Name.ValueString(), integrationTarget.ResourceType.ValueString()))
	return nil, diagnostics
}

func convertFiltersToListTagsApiModel(filters []models.ResourceFilter) ([]apono.TagV1, diag.Diagnostics) {
	data := []apono.TagV1{}
	for _, filter := range filters {
		switch filter.Type.ValueString() {
		case "tag":
			data = append(data, apono.TagV1{
				Name:  filter.Name.ValueString(),
				Value: filter.Value.ValueString(),
			})
		case "id":
			data = append(data, apono.TagV1{
				Name:  "__id",
				Value: filter.Value.ValueString(),
			})
		case "name":
			data = append(data, apono.TagV1{
				Name:  "__name",
				Value: filter.Value.ValueString(),
			})

		default:
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Filter Error", "Unsupported filter type, supported types are: tag, id, name")
			return nil, diagnostics
		}
	}

	return data, nil
}

func convertTagsToFiltersModel(tags []apono.TagV1) []models.ResourceFilter {
	var filters []models.ResourceFilter
	for _, tag := range tags {
		switch tag.Name {
		case "__id":
			filters = append(filters, models.ResourceFilter{
				Type:  types.StringValue("id"),
				Value: types.StringValue(tag.Value),
			})
		case "__name":
			filters = append(filters, models.ResourceFilter{
				Type:  types.StringValue("name"),
				Value: types.StringValue(tag.Value),
			})
		default:
			filters = append(filters, models.ResourceFilter{
				Type:  types.StringValue("tag"),
				Name:  types.StringValue(tag.Name),
				Value: types.StringValue(tag.Value),
			})
		}
	}

	return filters
}

func convertTriggerToApiModel(trigger models.Trigger) (*apono.AccessFlowTriggerV1, diag.Diagnostics) {
	var data apono.AccessFlowTriggerV1
	data.Type = trigger.Type.ValueString()

	if trigger.Timeframe != nil {
		var timeframeDays []apono.DayOfWeek
		for _, day := range trigger.Timeframe.DaysInWeek.Elements() {
			timeframeDays = append(timeframeDays, apono.DayOfWeek(utils.AttrValueToString(day)))
		}

		startOfDayTimeInSeconds, _ := trigger.Timeframe.StartOfDayTimeInSeconds.ValueBigFloat().Int64()
		endOfDayTimeInSeconds, _ := trigger.Timeframe.EndOfDayTimeInSeconds.ValueBigFloat().Int64()
		dataTimeFrame := apono.AccessFlowTriggerV1Timeframe{
			StartOfDayTimeInSeconds: startOfDayTimeInSeconds,
			EndOfDayTimeInSeconds:   endOfDayTimeInSeconds,
			DaysInWeek:              timeframeDays,
			TimeZone:                trigger.Timeframe.TimeZone.ValueString(),
		}

		data.Timeframe = *apono.NewNullableAccessFlowTriggerV1Timeframe(&dataTimeFrame)
	}

	return &data, nil
}

func getIdentitiesIdsByNameAndType(ctx context.Context, identityName string, identityType string, availableIdentities []apono.IdentityModel2, aponoClient *apono.APIClient) ([]string, diag.Diagnostics) {
	switch identityType {
	case "user":
		user, _, err := aponoClient.UsersApi.GetUser(ctx, identityName).Execute()
		if err != nil {
			return nil, utils.GetDiagnosticsForApiError(err, "get", "user", identityName)
		}

		return []string{user.GetId()}, nil

	case "group", "context_attribute":
		var identitiesIds []string
		for _, identity := range availableIdentities {
			if identity.Name == identityName && strings.ToLower(identity.Type) == identityType {
				identitiesIds = append(identitiesIds, identity.Id)
			}
		}
		if len(identitiesIds) == 0 {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get identity %s: %s", identityType, identityName))
			return nil, diagnostics
		}

		return identitiesIds, nil

	default:
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unsupported indentity type: %s, please update the provider to support this type or removed it using UI/API", identityType),
		)
		return nil, diagnostics
	}
}

func getUniqueListOfIdentities(identities []models.Identity) []models.Identity {
	var uniqueIdentities []models.Identity
	existingKeys := make(map[models.Identity]bool)
	for _, identity := range identities {
		if existingKeys[identity] {
			continue
		}
		uniqueIdentities = append(uniqueIdentities, identity)
		existingKeys[identity] = true
	}

	return uniqueIdentities
}
