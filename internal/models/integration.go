package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
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
