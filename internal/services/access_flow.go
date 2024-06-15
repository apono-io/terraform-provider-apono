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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"math/big"
	"strings"
)

func ConvertAccessFlowApiToTerraformModel(ctx context.Context, aponoClient *apono.APIClient, accessFlow *aponoapi.AccessFlowTerraformV1) (*models.AccessFlowModel, diag.Diagnostics) {
	revokeAfterInSec := types.NumberValue(big.NewFloat(float64(accessFlow.GetRevokeAfterInSec())))

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

		startTimeString := utils.SecondsToDayTimeFormat(int(timeframe.GetStartOfDayTimeInSeconds()))
		endTimeString := utils.SecondsToDayTimeFormat(int(timeframe.GetEndOfDayTimeInSeconds()))

		dataTrigger.Timeframe = &models.Timeframe{
			StartTime:  types.StringValue(startTimeString),
			EndTime:    types.StringValue(endTimeString),
			DaysInWeek: daysInWeek,
			TimeZone:   types.StringValue(timeframe.GetTimeZone()),
		}
	}

	availableIdentities, _, err := aponoClient.IdentitiesApi.ListIdentities(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "identities", "")
	}

	availableUsers, _, err := aponoClient.UsersApi.ListUsers(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "users", "")
	}

	var dataGrantees []models.Identity
	for _, grantee := range accessFlow.GetGrantees() {
		identity, diagnostics := convertIdentityApiToTerraformModel(grantee.Id, strings.ToLower(grantee.Type), availableIdentities.Data, availableUsers.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataGrantees = append(dataGrantees, *identity)
	}

	// This converts the list of identities to a Terraform Set, which require map of attribute name to type.
	setGrantees, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: models.IdentityObject}, getUniqueListOfIdentities(dataGrantees))
	if len(diags) > 0 {
		return nil, diags
	}

	dataGranteeFilterGroup, diags := convertGranteeFilterGroupApiToTerraformModel(ctx, accessFlow.GranteeFilterGroup.Get())
	if len(diags) > 0 {
		return nil, diags
	}
	objectGranteeFilterGroup, diags := types.ObjectValueFrom(ctx, models.GranteeFilterGroupObject, dataGranteeFilterGroup)
	if len(diags) > 0 {
		return nil, diags
	}

	var dataApprovers []models.Identity
	for _, approver := range accessFlow.GetApprovers() {
		identity, diagnostics := convertIdentityApiToTerraformModel(approver.Id, strings.ToLower(approver.Type), availableIdentities.Data, availableUsers.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataApprovers = append(dataApprovers, *identity)
	}

	// This converts the list of identities to a Terraform Set, which require map of attribute name to type.
	setApprovers, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: models.IdentityObject}, getUniqueListOfIdentities(dataApprovers))
	if len(diags) > 0 {
		return nil, diags
	}

	apiIntegrationTargets := convertIntegrationTargetsNewApiToOldApiModel(accessFlow.GetIntegrationTargets())
	dataIntegrationTargets, diagnostics := convertIntegrationTargetsApiToTerraformModel(ctx, aponoClient, apiIntegrationTargets)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	dataBundleTargets, diagnostics := convertBundleTargetsApiToTerraformModel(ctx, aponoClient, accessFlow.GetBundleTargets())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var dataSettings *models.Settings
	existingSettings, ok := accessFlow.GetSettingsOk()
	if ok && existingSettings != nil {
		dataSettings = &models.Settings{
			RequireJustificationOnRequestAgain: types.BoolValue(existingSettings.GetRequireJustificationOnRequestAgain()),
			RequireAllApprovers:                types.BoolValue(existingSettings.GetRequireAllApprovers()),
			ApproverCannotSelfApprove:          types.BoolValue(existingSettings.GetApproverCannotApproveHimself()),
		}
	} else {
		dataSettings = nil
	}

	dataLabels, diagnostics := convertLabelsApiToTerraformModel(ctx, accessFlow.GetLabels())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := models.AccessFlowModel{
		ID:                  types.StringValue(accessFlow.GetId()),
		Name:                types.StringValue(accessFlow.GetName()),
		Active:              types.BoolValue(accessFlow.GetActive()),
		RevokeAfterInSec:    revokeAfterInSec,
		Trigger:             &dataTrigger,
		Grantees:            setGrantees,
		GranteesFilterGroup: objectGranteeFilterGroup,
		IntegrationTargets:  dataIntegrationTargets,
		BundleTargets:       dataBundleTargets,
		Approvers:           setApprovers,
		Settings:            dataSettings,
		Labels:              *dataLabels,
	}

	return &data, nil
}

func ConvertAccessFlowTerraformModelToApi(ctx context.Context, aponoClient *apono.APIClient, accessFlow *models.AccessFlowModel) (*aponoapi.UpsertAccessFlowTerraformV1, diag.Diagnostics) {
	revokeAfterInSec, _ := accessFlow.RevokeAfterInSec.ValueBigFloat().Int64()

	dataTrigger, diagnostics := convertTriggerTerraformModelToApi(*accessFlow.Trigger)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	availableIdentities, _, err := aponoClient.IdentitiesApi.ListIdentities(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "identities", "")
	}

	availableUsers, _, err := aponoClient.UsersApi.ListUsers(ctx).Execute()
	if err != nil {
		return nil, utils.GetDiagnosticsForApiError(err, "list", "users", "")
	}

	existingGrantees := make([]models.Identity, 0, len(accessFlow.Grantees.Elements()))
	if !accessFlow.Grantees.IsNull() && !accessFlow.Grantees.IsUnknown() {
		diagnostics = accessFlow.Grantees.ElementsAs(ctx, &existingGrantees, false)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
	}

	dataGrantees := []aponoapi.GranteeTerraformV1{}
	for _, grantee := range existingGrantees {
		granteeIds, diagnostics := getIdentitiesIdsByNameAndType(grantee.Name.ValueString(), grantee.Type.ValueString(), availableIdentities.Data, availableUsers.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		for _, granteeId := range granteeIds {
			dataGrantees = append(dataGrantees, aponoapi.GranteeTerraformV1{
				Id:   granteeId,
				Type: grantee.Type.ValueString(),
			})
		}
	}

	var dataGranteeFilterGroup *aponoapi.AccessFlowTerraformV1GranteeFilterGroup
	if !accessFlow.GranteesFilterGroup.IsNull() && !accessFlow.GranteesFilterGroup.IsUnknown() {
		var modelGranteeFilterGroup models.GranteeFilterGroup
		diagnostics = accessFlow.GranteesFilterGroup.As(ctx, &modelGranteeFilterGroup, basetypes.ObjectAsOptions{})
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
		dataGranteeFilterGroup, diagnostics = convertGranteeFilterGroupTerraformModelToApi(ctx, &modelGranteeFilterGroup)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}
	}

	existingApprovers := make([]models.Identity, 0, len(accessFlow.Approvers.Elements()))
	diagnostics = accessFlow.Approvers.ElementsAs(ctx, &existingApprovers, false)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var dataApprovers []aponoapi.ApproverTerraformV1
	for _, approver := range existingApprovers {
		approverIds, diagnostics := getIdentitiesIdsByNameAndType(approver.Name.ValueString(), approver.Type.ValueString(), availableIdentities.Data, availableUsers.Data)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		for _, approverId := range approverIds {
			dataApprovers = append(dataApprovers, aponoapi.ApproverTerraformV1{
				Id:   approverId,
				Type: approver.Type.ValueString(),
			})
		}
	}

	dataIntegrationTargets, diagnostics := convertIntegrationTargetsTerraformModelToApi(ctx, aponoClient, accessFlow.IntegrationTargets)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	dataBundleTargets, diagnostics := convertBundleTargetsTerraformModelToApi(ctx, aponoClient, accessFlow.BundleTargets)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	setting := aponoapi.NullableAccessFlowTerraformV1Settings{}
	if accessFlow.Settings != nil {
		settings := aponoapi.AccessFlowTerraformV1Settings{
			RequireJustificationOnRequestAgain: *aponoapi.NewNullableBool(accessFlow.Settings.RequireJustificationOnRequestAgain.ValueBoolPointer()),
			RequireAllApprovers:                *aponoapi.NewNullableBool(accessFlow.Settings.RequireAllApprovers.ValueBoolPointer()),
			ApproverCannotApproveHimself:       *aponoapi.NewNullableBool(accessFlow.Settings.ApproverCannotSelfApprove.ValueBoolPointer()),
		}

		setting.Set(&settings)
	} else {
		setting.Unset()
	}

	dataLabels := convertLabelsTerraformModelToApi(accessFlow.Labels)

	data := aponoapi.UpsertAccessFlowTerraformV1{
		Name:               accessFlow.Name.ValueString(),
		Active:             accessFlow.Active.ValueBool(),
		RevokeAfterInSec:   int32(revokeAfterInSec),
		Trigger:            *dataTrigger,
		Grantees:           dataGrantees,
		GranteeFilterGroup: *aponoapi.NewNullableAccessFlowTerraformV1GranteeFilterGroup(dataGranteeFilterGroup),
		Approvers:          dataApprovers,
		IntegrationTargets: convertIntegrationTargetsOldApiToNewApiModel(dataIntegrationTargets),
		BundleTargets:      dataBundleTargets,
		Settings:           setting,
		Labels:             dataLabels,
	}

	return &data, nil
}

func convertIdentityApiToTerraformModel(identityId string, identityType string, availableIdentities []apono.IdentityModel2, availableUsers []apono.UserModel) (*models.Identity, diag.Diagnostics) {
	switch identityType {
	case "user":
		var result *models.Identity
		for _, user := range availableUsers {
			if user.Id == identityId {
				result = &models.Identity{
					Name: types.StringValue(user.GetEmail()),
					Type: types.StringValue("user"),
				}
			}
		}

		if result == nil {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get user: %s", identityId))
			return nil, diagnostics
		}

		return result, nil

	case "group", "context_attribute":
		var result *models.Identity
		for _, identity := range availableIdentities {
			if identity.Id == identityId {
				result = &models.Identity{
					Name: types.StringValue(identity.GetName()),
					Type: types.StringValue(identityType),
				}
			}
		}

		if result == nil {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get identity %s: %s", identityType, identityId))
			return nil, diagnostics
		}

		return result, nil

	default:
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unsupported indentity type: %s, please update the provider to support this type or removed it using UI/API", identityType),
		)
		return nil, diagnostics
	}
}

func convertResourceFilterListToTagV1Api(filters []models.ResourceFilter) ([]apono.TagV1, diag.Diagnostics) {
	data := []apono.TagV1{}
	for _, filter := range filters {
		switch filter.Type.ValueString() {
		case "tag":
			data = append(data, apono.TagV1{
				Name:  filter.Key.ValueString(),
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

func convertTagV1ListToResourceFilter(tags []apono.TagV1) []models.ResourceFilter {
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
				Key:   types.StringValue(tag.Name),
				Value: types.StringValue(tag.Value),
			})
		}
	}

	return filters
}

func convertTriggerTerraformModelToApi(trigger models.Trigger) (*aponoapi.AccessFlowTriggerTerraformV1, diag.Diagnostics) {
	var data aponoapi.AccessFlowTriggerTerraformV1
	data.Type = trigger.Type.ValueString()

	if trigger.Timeframe != nil {
		var timeframeDays []aponoapi.DayOfWeek
		for _, day := range trigger.Timeframe.DaysInWeek.Elements() {
			timeframeDays = append(timeframeDays, aponoapi.DayOfWeek(utils.AttrValueToString(day)))
		}

		startOfDayTimeInSeconds, err := utils.DayTimeFormatToSeconds(trigger.Timeframe.StartTime.ValueString())
		if err != nil {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to parse start time: %s", trigger.Timeframe.StartTime.ValueString()))
			return nil, diagnostics
		}
		endOfDayTimeInSeconds, err := utils.DayTimeFormatToSeconds(trigger.Timeframe.EndTime.ValueString())
		if err != nil {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to parse end time: %s", trigger.Timeframe.EndTime.ValueString()))
			return nil, diagnostics
		}

		dataTimeFrame := aponoapi.AccessFlowTriggerTerraformV1Timeframe{
			StartOfDayTimeInSeconds: startOfDayTimeInSeconds,
			EndOfDayTimeInSeconds:   endOfDayTimeInSeconds,
			DaysInWeek:              timeframeDays,
			TimeZone:                trigger.Timeframe.TimeZone.ValueString(),
		}

		data.Timeframe = *aponoapi.NewNullableAccessFlowTriggerTerraformV1Timeframe(&dataTimeFrame)
	}

	return &data, nil
}

func getIdentitiesIdsByNameAndType(identityName string, identityType string, availableIdentities []apono.IdentityModel2, availableUsers []apono.UserModel) ([]string, diag.Diagnostics) {
	switch identityType {
	case "user":
		var userId string
		for _, user := range availableUsers {
			if user.Email == identityName {
				userId = user.GetId()
			}
		}

		if userId == "" {
			diagnostics := diag.Diagnostics{}
			diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get user: %s", identityName))
			return nil, diagnostics
		}

		return []string{userId}, nil

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

// getUniqueListOfIdentities returns a unique list of identities.
// This is used because the API returns duplicate identities in case of groups with the same name.
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

func convertGranteeFilterGroupApiToTerraformModel(ctx context.Context, granteeFilterGroup *aponoapi.AccessFlowTerraformV1GranteeFilterGroup) (*models.GranteeFilterGroup, diag.Diagnostics) {
	if granteeFilterGroup == nil {
		return nil, nil
	}

	var dataFilters []models.AttributeFilter
	for _, apiFilter := range granteeFilterGroup.GetAttributeFilters() {
		dataFilter, diagnostics := convertAttributeFiltersTerraformModelToApi(ctx, apiFilter)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		dataFilters = append(dataFilters, *dataFilter)
	}

	filters, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: models.AttributeFilterObject}, dataFilters)
	if len(diags) > 0 {
		return nil, diags
	}

	return &models.GranteeFilterGroup{
		Operator: types.StringValue(string(granteeFilterGroup.GetLogicalOperator())),
		Filters:  filters,
	}, nil
}

func convertGranteeFilterGroupTerraformModelToApi(ctx context.Context, granteeFilterGroup *models.GranteeFilterGroup) (*aponoapi.AccessFlowTerraformV1GranteeFilterGroup, diag.Diagnostics) {
	if granteeFilterGroup == nil {
		return nil, nil
	}

	dataFilters := make([]models.AttributeFilter, 0, len(granteeFilterGroup.Filters.Elements()))
	diagnostics := granteeFilterGroup.Filters.ElementsAs(ctx, &dataFilters, false)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var apiFilters []aponoapi.AttributeFilterTerraformV1
	for _, dataFilter := range dataFilters {
		apiFilter, diagnostics := convertAttributeFiltersApiToTerraformModel(ctx, dataFilter)
		if len(diagnostics) > 0 {
			return nil, diagnostics
		}

		apiFilters = append(apiFilters, *apiFilter)
	}

	return &aponoapi.AccessFlowTerraformV1GranteeFilterGroup{
		LogicalOperator:  aponoapi.GranteeFilterGroupOperatorTerraformV1(granteeFilterGroup.Operator.ValueString()),
		AttributeFilters: apiFilters,
	}, nil
}

func convertAttributeFiltersApiToTerraformModel(ctx context.Context, filter models.AttributeFilter) (*aponoapi.AttributeFilterTerraformV1, diag.Diagnostics) {
	attributeValues := make([]string, 0, len(filter.AttributeNames.Elements()))
	diagnostics := filter.AttributeNames.ElementsAs(ctx, &attributeValues, false)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	var operator *string
	if !filter.Operator.IsNull() && !filter.Operator.IsUnknown() {
		operator = filter.Operator.ValueStringPointer()
	}

	return &aponoapi.AttributeFilterTerraformV1{
		Operator:        *aponoapi.NewNullableString(operator),
		AttributeTypeId: filter.AttributeType.ValueString(),
		AttributeValue:  attributeValues,
		IntegrationId:   *aponoapi.NewNullableString(filter.IntegrationID.ValueStringPointer()),
	}, nil
}

func convertAttributeFiltersTerraformModelToApi(ctx context.Context, filter aponoapi.AttributeFilterTerraformV1) (*models.AttributeFilter, diag.Diagnostics) {
	apiFilterNames, err := utils.ConvertInterfaceToListOfString(filter.GetAttributeValue())
	if err != nil {
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Client Error", fmt.Sprintf("Failed to convert attribute names for attribute filter: %s", err.Error()))
		return nil, diagnostics
	}
	dataFilterNames, diagnostics := types.ListValueFrom(ctx, types.StringType, apiFilterNames)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	return &models.AttributeFilter{
		Operator:       types.StringPointerValue(filter.Operator.Get()),
		AttributeType:  types.StringValue(filter.AttributeTypeId),
		AttributeNames: dataFilterNames,
		IntegrationID:  types.StringPointerValue(filter.IntegrationId.Get()),
	}, nil
}

func convertLabelsApiToTerraformModel(ctx context.Context, labels []aponoapi.AccessFlowLabelTerraformV1) (*basetypes.ListValue, diag.Diagnostics) {
	var labelNames []string
	for _, label := range labels {
		labelNames = append(labelNames, label.Name)
	}

	if len(labelNames) == 0 {
		undefinedList := basetypes.NewListNull(types.StringType)
		return &undefinedList, nil
	}

	dataLabels, diagnostics := types.ListValueFrom(ctx, types.StringType, labelNames)
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	return &dataLabels, nil
}

func convertLabelsTerraformModelToApi(labels basetypes.ListValue) []aponoapi.AccessFlowLabelTerraformV1 {
	data := []aponoapi.AccessFlowLabelTerraformV1{}
	for _, label := range labels.Elements() {
		data = append(data, aponoapi.AccessFlowLabelTerraformV1{
			Name: utils.AttrValueToString(label),
		})
	}

	return data
}
