package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/semaphoreui/semaphore/db"
)

type AzureAccessKeyDeserializer struct {
	accessKeyRepo     db.AccessKeyManager
	secretStorageRepo db.SecretStorageRepository
	encryptionService AccessKeyEncryptionService
}

func NewAzureAccessKeyDeserializer(
	accessKeyRepo db.AccessKeyManager,
	secretStorageRepo db.SecretStorageRepository,
	encryptionService AccessKeyEncryptionService,
) *AzureAccessKeyDeserializer {
	return &AzureAccessKeyDeserializer{
		accessKeyRepo:     accessKeyRepo,
		secretStorageRepo: secretStorageRepo,
		encryptionService: encryptionService,
	}
}

func (d *AzureAccessKeyDeserializer) getAzureClient(storage *db.SecretStorage) (*azsecrets.Client, error) {
	var params db.AzureSecretStorageParams
	err := storage.ExtractParams(&params)
	if err != nil {
		return nil, err
	}

	if params.VaultURL == "" {
		return nil, fmt.Errorf("vault URL is required for Azure Key Vault")
	}

	if params.TenantID == "" || params.ClientID == "" || params.ClientSecret == "" {
		return nil, fmt.Errorf("tenant ID, client ID, and client secret are required for Azure Key Vault")
	}

	// Create credentials using client secret
	cred, err := azidentity.NewClientSecretCredential(
		params.TenantID,
		params.ClientID,
		params.ClientSecret,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	// Create Azure Key Vault client
	client, err := azsecrets.NewClient(params.VaultURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Key Vault client: %w", err)
	}

	return client, nil
}

func (d *AzureAccessKeyDeserializer) getSecretName(key *db.AccessKey) string {
	if key.SourceStorageKey != nil && *key.SourceStorageKey != "" {
		return *key.SourceStorageKey
	}
	return fmt.Sprintf("semaphore-key-%d", key.ID)
}

func (d *AzureAccessKeyDeserializer) DeleteSecret(key *db.AccessKey) error {
	if key.SourceStorageID == nil {
		return fmt.Errorf("source storage ID is required for Azure Key Vault")
	}

	storage, err := d.secretStorageRepo.GetSecretStorage(*key.ProjectID, *key.SourceStorageID)
	if err != nil {
		return err
	}

	client, err := d.getAzureClient(&storage)
	if err != nil {
		return err
	}

	secretName := d.getSecretName(key)
	
	// Start deletion process (Azure Key Vault soft-deletes by default)
	_, err = client.DeleteSecret(context.Background(), secretName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete secret from Azure Key Vault: %w", err)
	}

	return nil
}

func (d *AzureAccessKeyDeserializer) SerializeSecret(key *db.AccessKey) error {
	if key.SourceStorageID == nil {
		return fmt.Errorf("source storage ID is required for Azure Key Vault")
	}

	storage, err := d.secretStorageRepo.GetSecretStorage(*key.ProjectID, *key.SourceStorageID)
	if err != nil {
		return err
	}

	client, err := d.getAzureClient(&storage)
	if err != nil {
		return err
	}

	var plaintext []byte

	switch key.Type {
	case db.AccessKeyString:
		if key.String == "" {
			// If no string value, delete the secret and clear the reference
			if err := d.DeleteSecret(key); err != nil {
				return err
			}
			key.Secret = nil
			return nil
		}
		plaintext = []byte(key.String)
	case db.AccessKeySSH:
		if key.SshKey.PrivateKey == "" {
			if key.SshKey.Login != "" || key.SshKey.Passphrase != "" {
				return fmt.Errorf("invalid ssh key")
			}
			// If no private key, delete the secret and clear the reference
			if err := d.DeleteSecret(key); err != nil {
				return err
			}
			key.Secret = nil
			return nil
		}

		plaintext, err = json.Marshal(key.SshKey)
		if err != nil {
			return err
		}
	case db.AccessKeyLoginPassword:
		if key.LoginPassword.Password == "" {
			if key.LoginPassword.Login != "" {
				return fmt.Errorf("invalid password key")
			}
			// If no password, delete the secret and clear the reference
			if err := d.DeleteSecret(key); err != nil {
				return err
			}
			key.Secret = nil
			return nil
		}

		plaintext, err = json.Marshal(key.LoginPassword)
		if err != nil {
			return err
		}
	case db.AccessKeyNone:
		key.Secret = nil
		return nil
	default:
		return fmt.Errorf("invalid access token type")
	}

	secretName := d.getSecretName(key)
	secretValue := string(plaintext)

	// Store secret in Azure Key Vault
	_, err = client.SetSecret(context.Background(), secretName, azsecrets.SetSecretParameters{
		Value: &secretValue,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to store secret in Azure Key Vault: %w", err)
	}

	// Set the secret reference in the key
	key.Secret = &secretName

	return nil
}

func (d *AzureAccessKeyDeserializer) DeserializeSecret(key *db.AccessKey) (string, error) {
	if key.Secret == nil || *key.Secret == "" {
		return "", nil
	}

	if key.SourceStorageID == nil {
		return "", fmt.Errorf("source storage ID is required for Azure Key Vault")
	}

	storage, err := d.secretStorageRepo.GetSecretStorage(*key.ProjectID, *key.SourceStorageID)
	if err != nil {
		return "", err
	}

	client, err := d.getAzureClient(&storage)
	if err != nil {
		return "", err
	}

	secretName := *key.Secret

	// Retrieve secret from Azure Key Vault
	resp, err := client.GetSecret(context.Background(), secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret from Azure Key Vault: %w", err)
	}

	if resp.Value == nil {
		return "", fmt.Errorf("secret value is nil")
	}

	return *resp.Value, nil
}