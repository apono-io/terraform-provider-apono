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
	HashicorpVaultSecret   *HashicorpVaultSecret  `tfsdk:"hashicorp_vault_secret"`
	AzureSecret            *AzureSecret           `tfsdk:"azure_secret"`
	AponoSecret            *AponoSecret           `tfsdk:"apono_secret"`
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

type HashicorpVaultSecret struct {
	SecretEngine types.String `tfsdk:"secret_engine"`
	Path         types.String `tfsdk:"path"`
}

type AzureSecret struct {
	VaultURL types.String `tfsdk:"vault_url"`
	Name     types.String `tfsdk:"name"`
}

type AponoSecret struct {
	Params types.Map `tfsdk:"params"`
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
