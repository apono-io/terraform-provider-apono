package models

import (
	"context"
	"fmt"
	"github.com/apono-io/terraform-provider-apono/internal/aponoapi"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// IntegrationModel describes the resource data model.
type IntegrationModel struct {
	ID                     types.String           `tfsdk:"id"`
	Name                   types.String           `tfsdk:"name"`
	Type                   types.String           `tfsdk:"type"`
	ConnectorID            types.String           `tfsdk:"connector_id"`
	ConnectedResourceTypes types.Set              `tfsdk:"connected_resource_types"`
	Metadata               types.Map              `tfsdk:"metadata"`
	CustomAccessDetails    types.String           `tfsdk:"custom_access_details"`
	AwsSecret              *AwsSecret             `tfsdk:"aws_secret"`
	GcpSecret              *GcpSecret             `tfsdk:"gcp_secret"`
	KubernetesSecret       *KubernetesSecret      `tfsdk:"kubernetes_secret"`
	ResourceOwnerMappings  []ResourceOwnerMapping `tfsdk:"resource_owner_mappings"`
	IntegrationOwners      []IntegrationOwner     `tfsdk:"integration_owners"`
}

type AwsSecret struct {
	Region   types.String `tfsdk:"region"`
	SecretID types.String `tfsdk:"secret_id"`
}

type GcpSecret struct {
	Project  types.String `tfsdk:"project"`
	SecretID types.String `tfsdk:"secret_id"`
}

type KubernetesSecret struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
}

type ResourceOwnerMapping struct {
	TagName                types.String `tfsdk:"key_name"`
	AttributeType          types.String `tfsdk:"attribute"`
	AttributeIntegrationId types.String `tfsdk:"attribute_integration_id"`
}

type IntegrationOwner struct {
	IntegrationId  types.String   `tfsdk:"integration_id"`
	AttributeType  types.String   `tfsdk:"attribute"`
	AttributeValue []types.String `tfsdk:"value"`
}

func ConvertToIntegrationModel(ctx context.Context, integration *aponoapi.IntegrationTerraform) (*IntegrationModel, diag.Diagnostics) {
	metadataMapValue, diagnostics := types.MapValueFrom(ctx, types.StringType, integration.GetParams())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := IntegrationModel{}
	data.ID = types.StringValue(integration.GetId())
	data.Name = types.StringValue(integration.GetName())
	data.Type = types.StringValue(integration.GetType())
	data.ConnectorID = types.StringValue(integration.GetProvisionerId())
	data.Metadata = metadataMapValue

	if integration.CustomAccessDetails.IsSet() {
		data.CustomAccessDetails = types.StringValue(integration.GetCustomAccessDetails())
	}

	connectedResourceTypes, diagnostics := types.SetValueFrom(ctx, types.StringType, integration.GetConnectedResourceTypes())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}
	data.ConnectedResourceTypes = connectedResourceTypes

	secretConfig := integration.GetSecretConfig()
	switch secretConfig["type"] {
	case "AWS":
		data.AwsSecret = &AwsSecret{
			Region:   basetypes.NewStringValue(toString(secretConfig["region"])),
			SecretID: basetypes.NewStringValue(toString(secretConfig["secret_id"])),
		}
	case "GCP":
		data.GcpSecret = &GcpSecret{
			Project:  basetypes.NewStringValue(toString(secretConfig["project"])),
			SecretID: basetypes.NewStringValue(toString(secretConfig["secret_id"])),
		}
	case "KUBERNETES":
		data.KubernetesSecret = &KubernetesSecret{
			Namespace: basetypes.NewStringValue(toString(secretConfig["namespace"])),
			Name:      basetypes.NewStringValue(toString(secretConfig["name"])),
		}
	}

	data.ResourceOwnerMappings = ConvertResourceOwnersMappingToModel(integration.ResourceOwnersMappings)
	data.IntegrationOwners = ConvertIntegrationOwnerToData(integration.IntegrationOwners)

	return &data, nil
}

func ConvertIntegrationOwnerToData(Owners []aponoapi.IntegrationOwnerTerraform) []IntegrationOwner {
	if Owners == nil {
		return nil
	}
	var integrationOwners = make([]IntegrationOwner, len(Owners))
	for i, in := range Owners {
		var AttributeValue = make([]types.String, len(in.AttributeValue))
		for j, av := range in.AttributeValue {
			AttributeValue[j] = basetypes.NewStringValue(av)
		}
		var IntegrationId types.String
		if in.IntegrationId.IsSet() {
			IntegrationId = basetypes.NewStringValue(*in.IntegrationId.Get())
		} else {
			IntegrationId = basetypes.NewStringNull()
		}
		integrationOwners[i] = IntegrationOwner{
			IntegrationId:  IntegrationId,
			AttributeType:  basetypes.NewStringValue(in.AttributeType),
			AttributeValue: AttributeValue,
		}
	}
	return integrationOwners
}

func ConvertResourceOwnersMappingToModel(ResourceOwnersMappings []aponoapi.ResourceOwnerMappingTerraform) []ResourceOwnerMapping {
	if ResourceOwnersMappings == nil {
		return nil
	}
	result := make([]ResourceOwnerMapping, len(ResourceOwnersMappings))
	for i, r := range ResourceOwnersMappings {
		var AttributeIntegrationId types.String
		if r.AttributeIntegrationId.IsSet() {
			AttributeIntegrationId = basetypes.NewStringValue(*r.AttributeIntegrationId.Get())
		} else {
			AttributeIntegrationId = basetypes.NewStringNull()
		}
		result[i] = ResourceOwnerMapping{
			TagName:                basetypes.NewStringValue(r.TagName),
			AttributeType:          basetypes.NewStringValue(r.AttributeType),
			AttributeIntegrationId: AttributeIntegrationId,
		}
	}

	return result
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
