package server

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/semaphoreui/semaphore/db"
)

type AzureService interface {
	GetSubscriptions(config db.AzureConfig) ([]db.AzureSubscription, error)
	GetKeyVaults(config db.AzureConfig, subscriptionID string) ([]db.AzureKeyVault, error)
}

type azureService struct {
	azureConfigRepo db.AzureConfigManager
}

func NewAzureService(azureConfigRepo db.AzureConfigManager) AzureService {
	return &azureService{
		azureConfigRepo: azureConfigRepo,
	}
}

func (s *azureService) createCredential(config db.AzureConfig) (*azidentity.ClientSecretCredential, error) {
	return azidentity.NewClientSecretCredential(
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		nil,
	)
}

func (s *azureService) GetSubscriptions(config db.AzureConfig) ([]db.AzureSubscription, error) {
	cred, err := s.createCredential(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions client: %w", err)
	}

	var subscriptions []db.AzureSubscription
	pager := client.NewListPager(nil)
	
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions: %w", err)
		}
		
		for _, sub := range page.Value {
			if sub.SubscriptionID != nil && sub.DisplayName != nil && sub.State != nil {
				subscriptions = append(subscriptions, db.AzureSubscription{
					ID:            *sub.SubscriptionID,
					Name:          *sub.DisplayName,
					State:         string(*sub.State),
					TenantID:      config.TenantID,
					AzureConfigID: config.ID,
				})
			}
		}
	}

	return subscriptions, nil
}

func (s *azureService) GetKeyVaults(config db.AzureConfig, subscriptionID string) ([]db.AzureKeyVault, error) {
	cred, err := s.createCredential(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %w", err)
	}

	client, err := armkeyvault.NewVaultsClient(subscriptionID, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create key vault client: %w", err)
	}

	var keyVaults []db.AzureKeyVault
	pager := client.NewListPager(nil)
	
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list key vaults: %w", err)
		}
		
		for _, vault := range page.Value {
			if vault.ID != nil && vault.Name != nil && vault.Location != nil {
				// Construct vault URL from name - standard Azure Key Vault URL pattern
				vaultURL := fmt.Sprintf("https://%s.vault.azure.net/", *vault.Name)
				
				keyVaults = append(keyVaults, db.AzureKeyVault{
					ID:             *vault.ID,
					Name:           *vault.Name,
					VaultURL:       vaultURL,
					Location:       *vault.Location,
					SubscriptionID: subscriptionID,
					AzureConfigID:  config.ID,
				})
			}
		}
	}

	return keyVaults, nil
}