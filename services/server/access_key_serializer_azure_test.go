package server

import (
	"testing"
	"github.com/semaphoreui/semaphore/db"
	"github.com/stretchr/testify/assert"
)

func TestAzureAccessKeyDeserializer_getSecretName(t *testing.T) {
	deserializer := &AzureAccessKeyDeserializer{}
	
	// Test with custom source storage key
	key := &db.AccessKey{
		ID: 123,
		SourceStorageKey: stringPtr("custom-key-name"),
	}
	
	secretName := deserializer.getSecretName(key)
	assert.Equal(t, "custom-key-name", secretName)
	
	// Test with default naming (no source storage key)
	key.SourceStorageKey = nil
	secretName = deserializer.getSecretName(key)
	assert.Equal(t, "semaphore-key-123", secretName)
	
	// Test with empty source storage key
	key.SourceStorageKey = stringPtr("")
	secretName = deserializer.getSecretName(key)
	assert.Equal(t, "semaphore-key-123", secretName)
}

func TestAzureSecretStorageParams_ExtractParams(t *testing.T) {
	// Test Azure params extraction
	storage := &db.SecretStorage{
		Type: db.SecretStorageTypeAzure,
		Params: db.MapStringAnyField{
			"vault_url":     "https://test.vault.azure.net/",
			"tenant_id":     "test-tenant",
			"client_id":     "test-client",
			"client_secret": "test-secret",
		},
	}
	
	var params db.AzureSecretStorageParams
	err := storage.ExtractParams(&params)
	
	assert.NoError(t, err)
	assert.Equal(t, "https://test.vault.azure.net/", params.VaultURL)
	assert.Equal(t, "test-tenant", params.TenantID)
	assert.Equal(t, "test-client", params.ClientID)
	assert.Equal(t, "test-secret", params.ClientSecret)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}