package resources

import (
	"testing"

	"github.com/apono-io/terraform-provider-apono/internal/v2/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSecretStoreConfigPresentButEmpty(t *testing.T) {
	t.Run("nil block - not present, no error", func(t *testing.T) {
		assert.False(t, secretStoreConfigPresentButEmpty(nil))
	})

	t.Run("aws set - present and valid", func(t *testing.T) {
		cfg := &models.SecretStoreConfig{
			AWS: &models.AWSSecretConfig{
				Region:   types.StringValue("us-east-1"),
				SecretID: types.StringValue("my-secret"),
			},
		}
		assert.False(t, secretStoreConfigPresentButEmpty(cfg))
	})

	t.Run("gcp set - present and valid", func(t *testing.T) {
		cfg := &models.SecretStoreConfig{
			GCP: &models.GCPSecretConfig{
				Project:  types.StringValue("my-project"),
				SecretID: types.StringValue("my-secret"),
			},
		}
		assert.False(t, secretStoreConfigPresentButEmpty(cfg))
	})

	t.Run("kubernetes set - present and valid", func(t *testing.T) {
		cfg := &models.SecretStoreConfig{
			Kubernetes: &models.KubernetesSecretConfig{
				Namespace: types.StringValue("default"),
				Name:      types.StringValue("my-secret"),
			},
		}
		assert.False(t, secretStoreConfigPresentButEmpty(cfg))
	})

	t.Run("all sub-providers nil - present but empty, error", func(t *testing.T) {
		cfg := &models.SecretStoreConfig{}
		assert.True(t, secretStoreConfigPresentButEmpty(cfg))
	})
}
