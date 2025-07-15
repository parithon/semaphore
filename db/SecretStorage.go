package db

import "encoding/json"

type SecretStorageType string

const (
	SecretStorageTypeLocal SecretStorageType = "local"
	SecretStorageTypeVault SecretStorageType = "vault"
	SecretStorageTypeAzure SecretStorageType = "azure"
)

type SecretStorage struct {
	ID        int               `db:"id" json:"id"`
	ProjectID int               `db:"project_id" json:"project_id"`
	Name      string            `db:"name" json:"name"`
	Type      SecretStorageType `db:"type" json:"type"`
	Params    MapStringAnyField `db:"params" json:"params"`
	ReadOnly  bool              `db:"readonly" json:"readonly"`

	VaultToken string `db:"-" json:"vault_token,omitempty" backup:"-"`
}

type VaultSecretStorageParams struct {
	URL string `json:"url"`
}

type AzureSecretStorageParams struct {
	VaultURL   string `json:"vault_url"`
	TenantID   string `json:"tenant_id"`
	ClientID   string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (s *SecretStorage) ExtractParams(target any) (err error) {
	content, err := json.Marshal(s.Params)
	if err != nil {
		return
	}

	switch target.(type) {
	case *VaultSecretStorageParams:
	case *AzureSecretStorageParams:
	default:
		err = &ValidationError{"invalid target type for extracting secret storage params"}
		return
	}

	err = json.Unmarshal(content, target)
	return
}
