package resources

import (
	"context"

	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// secretStoreConfigNotEmpty enforces, at plan time, that when the user provides
// a secret_store_config block they pick exactly one inner provider. Combined
// with the existing Conflicting() validator (which forbids more than one),
// this gives an "exactly one when present" rule that mirrors the API contract.
//
// Absence of the block is intentionally allowed - integrations may legitimately
// have no external secret store.
type secretStoreConfigNotEmpty struct{}

func (v secretStoreConfigNotEmpty) Description(_ context.Context) string {
	return "When secret_store_config is set, exactly one of aws/gcp/azure/hashicorp_vault/kubernetes must be specified."
}

func (v secretStoreConfigNotEmpty) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v secretStoreConfigNotEmpty) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config models.ResourceIntegrationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !secretStoreConfigPresentButEmpty(config.SecretStoreConfig) {
		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("secret_store_config"),
		"Empty secret_store_config block",
		"When secret_store_config is provided, exactly one of aws, gcp, azure, hashicorp_vault, or kubernetes must be set. Remove the block, or specify a provider.",
	)
}

// secretStoreConfigPresentButEmpty returns true iff the user wrote a
// secret_store_config block but did not set any inner provider.
func secretStoreConfigPresentButEmpty(s *models.SecretStoreConfig) bool {
	if s == nil {
		return false
	}
	return s.AWS == nil && s.GCP == nil && s.Azure == nil && s.HashicorpVault == nil && s.Kubernetes == nil
}
