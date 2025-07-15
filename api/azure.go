package api

import (
	"net/http"

	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"github.com/semaphoreui/semaphore/services/server"
	log "github.com/sirupsen/logrus"
)

type AzureController struct {
	AzureService server.AzureService
}

func NewAzureController(azureService server.AzureService) *AzureController {
	return &AzureController{
		AzureService: azureService,
	}
}

func getAzureConfigs(w http.ResponseWriter, r *http.Request) {
	currentUser := helpers.GetFromContext(r, "user").(*db.User)
	if !currentUser.Admin {
		log.Warn(currentUser.Username + " is not permitted to view Azure configurations")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	configs, err := helpers.Store(r).GetAzureConfigs(db.RetrieveQueryParams{
		Filter: r.URL.Query().Get("s"),
	})

	if err != nil {
		panic(err)
	}

	// Remove sensitive information before sending to client
	for i := range configs {
		configs[i].ClientSecret = ""
	}

	helpers.WriteJSON(w, http.StatusOK, configs)
}

func addAzureConfig(w http.ResponseWriter, r *http.Request) {
	var config db.AzureConfig
	if !helpers.Bind(w, r, &config) {
		return
	}

	currentUser := helpers.GetFromContext(r, "user").(*db.User)
	if !currentUser.Admin {
		log.Warn(currentUser.Username + " is not permitted to create Azure configurations")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newConfig, err := helpers.Store(r).CreateAzureConfig(config)
	if err != nil {
		log.Warn("Azure config creation failed: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove sensitive information before sending to client
	newConfig.ClientSecret = ""
	helpers.WriteJSON(w, http.StatusCreated, newConfig)
}

func getAzureConfigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configID, err := helpers.GetIntParam("config_id", w, r)
		if err != nil {
			return
		}

		config, err := helpers.Store(r).GetAzureConfig(configID)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		currentUser := helpers.GetFromContext(r, "user").(*db.User)
		if !currentUser.Admin {
			log.Warn(currentUser.Username + " is not permitted to access Azure configurations")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = helpers.SetContextValue(r, "_azure_config", config)
		next.ServeHTTP(w, r)
	})
}

func updateAzureConfig(w http.ResponseWriter, r *http.Request) {
	targetConfig := helpers.GetFromContext(r, "_azure_config").(db.AzureConfig)
	currentUser := helpers.GetFromContext(r, "user").(*db.User)

	var config db.AzureConfig
	if !helpers.Bind(w, r, &config) {
		return
	}

	if !currentUser.Admin {
		log.Warn(currentUser.Username + " is not permitted to edit Azure configurations")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	config.ID = targetConfig.ID
	config.Created = targetConfig.Created
	
	if err := helpers.Store(r).UpdateAzureConfig(config); err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteAzureConfig(w http.ResponseWriter, r *http.Request) {
	config := helpers.GetFromContext(r, "_azure_config").(db.AzureConfig)
	currentUser := helpers.GetFromContext(r, "user").(*db.User)

	if !currentUser.Admin {
		log.Warn(currentUser.Username + " is not permitted to delete Azure configurations")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := helpers.Store(r).DeleteAzureConfig(config.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getAzureConfig(w http.ResponseWriter, r *http.Request) {
	config := helpers.GetFromContext(r, "_azure_config").(db.AzureConfig)
	
	// Remove sensitive information before sending to client
	config.ClientSecret = ""
	helpers.WriteJSON(w, http.StatusOK, config)
}

func (c *AzureController) GetAzureSubscriptions(w http.ResponseWriter, r *http.Request) {
	config := helpers.GetFromContext(r, "_azure_config").(db.AzureConfig)
	
	subscriptions, err := c.AzureService.GetSubscriptions(config)
	if err != nil {
		log.WithError(err).Error("Failed to get Azure subscriptions")
		helpers.WriteErrorStatus(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, subscriptions)
}

func (c *AzureController) GetAzureKeyVaults(w http.ResponseWriter, r *http.Request) {
	config := helpers.GetFromContext(r, "_azure_config").(db.AzureConfig)
	subscriptionID := r.URL.Query().Get("subscription_id")
	
	if subscriptionID == "" {
		helpers.WriteErrorStatus(w, "subscription_id parameter is required", http.StatusBadRequest)
		return
	}
	
	keyVaults, err := c.AzureService.GetKeyVaults(config, subscriptionID)
	if err != nil {
		log.WithError(err).Error("Failed to get Azure key vaults")
		helpers.WriteErrorStatus(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, keyVaults)
}