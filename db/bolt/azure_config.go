package bolt

import (
	"github.com/semaphoreui/semaphore/db"
)

// Stub implementations for Azure configuration management in BoltDB
// These are not implemented as BoltDB is mainly used for testing/demo purposes

func (d *BoltDb) GetAzureConfigs(params db.RetrieveQueryParams) ([]db.AzureConfig, error) {
	return []db.AzureConfig{}, nil
}

func (d *BoltDb) CreateAzureConfig(config db.AzureConfig) (db.AzureConfig, error) {
	return db.AzureConfig{}, db.ErrNotFound
}

func (d *BoltDb) DeleteAzureConfig(configID int) error {
	return db.ErrNotFound
}

func (d *BoltDb) UpdateAzureConfig(config db.AzureConfig) error {
	return db.ErrNotFound
}

func (d *BoltDb) GetAzureConfig(configID int) (db.AzureConfig, error) {
	return db.AzureConfig{}, db.ErrNotFound
}