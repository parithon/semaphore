package sql

import (
	"github.com/Masterminds/squirrel"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/pkg/tz"
)

func (d *SqlDb) CreateAzureConfig(config db.AzureConfig) (newConfig db.AzureConfig, err error) {
	err = db.ValidateAzureConfig(config)
	if err != nil {
		return
	}

	config.Created = db.GetParsedTime(tz.Now())

	err = d.sql.Insert(&config)
	if err != nil {
		return
	}

	newConfig = config
	return
}

func (d *SqlDb) DeleteAzureConfig(configID int) error {
	res, err := d.exec("delete from `azure_config` where id=?", configID)
	return validateMutationResult(res, err)
}

func (d *SqlDb) UpdateAzureConfig(config db.AzureConfig) error {
	_, err := d.exec(
		"update `azure_config` set name=?, description=?, tenant_id=?, client_id=?, client_secret=? where id=?",
		config.Name,
		config.Description,
		config.TenantID,
		config.ClientID,
		config.ClientSecret,
		config.ID)

	return err
}

func (d *SqlDb) GetAzureConfig(configID int) (config db.AzureConfig, err error) {
	err = d.selectOne(&config, "select * from `azure_config` where id=?", configID)
	return
}

func (d *SqlDb) GetAzureConfigs(params db.RetrieveQueryParams) (configs []db.AzureConfig, err error) {
	q := squirrel.Select("*").From("`azure_config`")

	q, err = getQueryForParams(q, "", db.AzureConfigProps, params)
	if err != nil {
		return
	}

	if params.Filter != "" {
		q = q.Where(squirrel.Like{"name": escapeLike(params.Filter) + "%"})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return
	}

	_, err = d.selectAll(&configs, query, args...)
	return
}