package db

import (
	"time"
)

// AzureConfig represents an Azure configuration for connecting to Azure services
type AzureConfig struct {
	ID           int       `db:"id" json:"id"`
	Created      time.Time `db:"created" json:"created"`
	Name         string    `db:"name" json:"name" binding:"required"`
	Description  string    `db:"description" json:"description"`
	TenantID     string    `db:"tenant_id" json:"tenant_id" binding:"required"`
	ClientID     string    `db:"client_id" json:"client_id" binding:"required"`
	ClientSecret string    `db:"client_secret" json:"client_secret" binding:"required"`
}

// AzureSubscription represents an Azure subscription discovered from Azure API
type AzureSubscription struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	State        string `json:"state"`
	TenantID     string `json:"tenant_id"`
	AzureConfigID int   `json:"azure_config_id"`
}

// AzureKeyVault represents an Azure Key Vault discovered from Azure API
type AzureKeyVault struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	VaultURL       string `json:"vault_url"`
	Location       string `json:"location"`
	SubscriptionID string `json:"subscription_id"`
	AzureConfigID  int    `json:"azure_config_id"`
}

func ValidateAzureConfig(config AzureConfig) error {
	if config.Name == "" {
		return &ValidationError{Message: "Name cannot be empty"}
	}
	if config.TenantID == "" {
		return &ValidationError{Message: "Tenant ID cannot be empty"}
	}
	if config.ClientID == "" {
		return &ValidationError{Message: "Client ID cannot be empty"}
	}
	if config.ClientSecret == "" {
		return &ValidationError{Message: "Client Secret cannot be empty"}
	}
	return nil
}