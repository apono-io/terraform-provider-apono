package models

import (
	"context"
	"fmt"
	"github.com/apono-io/apono-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// IntegrationModel describes the resource data model.
type IntegrationModel struct {
	ID                     types.String      `tfsdk:"id"`
	Name                   types.String      `tfsdk:"name"`
	Type                   types.String      `tfsdk:"type"`
	ConnectorID            types.String      `tfsdk:"connector_id"`
	ConnectedResourceTypes types.Set         `tfsdk:"connected_resource_types"`
	Metadata               types.Map         `tfsdk:"metadata"`
	CustomAccessDetails    types.String      `tfsdk:"custom_access_details"`
	AwsSecret              *AwsSecret        `tfsdk:"aws_secret"`
	GcpSecret              *GcpSecret        `tfsdk:"gcp_secret"`
	KubernetesSecret       *KubernetesSecret `tfsdk:"kubernetes_secret"`
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

func ConvertToIntegrationModel(ctx context.Context, integration *apono.Integration) (*IntegrationModel, diag.Diagnostics) {
	metadataMapValue, diagnostics := types.MapValueFrom(ctx, types.StringType, integration.GetMetadata())
	if len(diagnostics) > 0 {
		return nil, diagnostics
	}

	data := IntegrationModel{}
	data.ID = types.StringValue(integration.GetId())
	data.Name = types.StringValue(integration.GetName())
	data.Type = types.StringValue(integration.GetType())
	data.ConnectorID = types.StringValue(integration.GetProvisionerId())
	data.Metadata = metadataMapValue

	if integration.CustomInstructionMessage.IsSet() {
		data.CustomAccessDetails = types.StringValue(integration.GetCustomInstructionMessage())
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

	return &data, nil
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
