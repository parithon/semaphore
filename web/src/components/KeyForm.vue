<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="item != null && secretStorages != null"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-text-field
        v-model="item.name"
        :label="$t('keyName')"
        :rules="[v => !!v || $t('name_required')]"
        required
        :disabled="formSaving"
        outlined
        dense
    />

    <v-autocomplete
      v-model="item.source_storage_id"
      :label="$t('Storage (optional)')"
      :items="secretStorages"
      item-value="id"
      item-text="name"
      :disabled="formSaving || !canEditSecrets"
      outlined
      dense
      clearable
    />

    <v-text-field
      v-if="item.source_storage_id != null"
      v-model="item.source_storage_key"
      :label="$t('Source Key')"
      :disabled="formSaving || !canEditSecrets"
      outlined
      dense
    />

    <!-- Azure-specific fields -->
    <div v-if="isAzureStorage && canEditSecrets">
      <v-autocomplete
        v-model="selectedAzureConfig"
        label="Azure Configuration"
        :items="azureConfigs"
        item-value="id"
        item-text="name"
        :disabled="formSaving"
        outlined
        dense
        clearable
        @input="onAzureConfigChange"
      />

      <v-autocomplete
        v-if="selectedAzureConfig"
        v-model="selectedSubscription"
        label="Azure Subscription"
        :items="azureSubscriptions"
        item-value="id"
        item-text="name"
        :disabled="formSaving || loadingSubscriptions"
        :loading="loadingSubscriptions"
        outlined
        dense
        clearable
        @input="onSubscriptionChange"
      />

      <v-autocomplete
        v-if="selectedSubscription"
        v-model="selectedKeyVault"
        label="Azure Key Vault"
        :items="azureKeyVaults"
        item-value="vault_url"
        item-text="name"
        :disabled="formSaving || loadingKeyVaults"
        :loading="loadingKeyVaults"
        outlined
        dense
        clearable
        @input="onKeyVaultChange"
      />

      <v-alert
        v-if="isAzureStorage && !selectedAzureConfig"
        dense
        text
        type="info"
        class="mt-2"
      >
        Please configure Azure settings in the Azure page first, then select an Azure configuration above.
      </v-alert>
    </div>

    <v-select
        v-model="item.type"
        :label="$t('type')"
        :rules="[v => (!!v || !canEditSecrets) || $t('type_required')]"
        :items="inventoryTypes"
        item-value="id"
        item-text="name"
        :required="canEditSecrets"
        :disabled="formSaving || !canEditSecrets"
        outlined
        dense
    />

    <v-text-field
        v-model="item.login_password.login"
        :label="$t('loginOptional')"
        v-if="item.type === 'login_password'"
        :disabled="formSaving || !canEditSecrets"
        outlined
        dense
    />

    <v-text-field
        v-model="item.login_password.password"
        :append-icon="showLoginPassword ? 'mdi-eye' : 'mdi-eye-off'"
        :label="$t('password')"
        :rules="[v => (!!v || !canEditSecrets) || $t('password_required')]"
        :type="showLoginPassword ? 'text' : 'password'"
        v-if="item.type === 'login_password'"
        :required="canEditSecrets"
        :disabled="formSaving || !canEditSecrets"
        autocomplete="new-password"
        @click:append="showLoginPassword = !showLoginPassword"
        outlined
        dense
    />

    <v-text-field
      v-model="item.ssh.login"
      :label="$t('usernameOptional')"
      v-if="item.type === 'ssh'"
      :disabled="formSaving || !canEditSecrets"
      outlined
      dense
    />

    <v-text-field
      v-model="item.ssh.passphrase"
      :append-icon="showSSHPassphrase ? 'mdi-eye' : 'mdi-eye-off'"
      label="Passphrase (Optional)"
      :type="showSSHPassphrase ? 'text' : 'password'"
      v-if="item.type === 'ssh'"
      :disabled="formSaving || !canEditSecrets"
      @click:append="showSSHPassphrase = !showSSHPassphrase"
      outlined
      dense
    />

    <v-textarea
      outlined
      v-model="item.ssh.private_key"
      :label="$t('privateKey')"
      :disabled="formSaving || !canEditSecrets"
      :rules="[v => !canEditSecrets || !!v || $t('private_key_required')]"
      v-if="item.type === 'ssh'"
    />

    <v-checkbox
        v-model="item.override_secret"
        :label="$t('override')"
        v-if="!isNew"
    />

    <v-alert
        dense
        text
        type="info"
        v-if="item.type === 'none'"
    >
      {{ $t('useThisTypeOfKeyForHttpsRepositoriesAndForPlaybook') }}
    </v-alert>
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],

  props: {
    supportStorages: Boolean,
  },

  data() {
    return {
      showLoginPassword: false,
      showSSHPassphrase: false,
      inventoryTypes: [{
        id: 'ssh',
        name: `${this.$t('keyFormSshKey')}`,
      }, {
        id: 'login_password',
        name: `${this.$t('keyFormLoginPassword')}`,
      }, {
        id: 'none',
        name: `${this.$t('keyFormNone')}`,
      }],
      secretStorages: null,
      azureConfigs: [],
      azureSubscriptions: [],
      azureKeyVaults: [],
      selectedAzureConfig: null,
      selectedSubscription: null,
      selectedKeyVault: null,
      loadingSubscriptions: false,
      loadingKeyVaults: false,
    };
  },

  computed: {
    canEditSecrets() {
      return this.isNew || this.item.override_secret;
    },
    
    isAzureStorage() {
      if (!this.item.source_storage_id || !this.secretStorages) return false;
      const storage = this.secretStorages.find(s => s.id === this.item.source_storage_id);
      return storage && storage.type === 'azure';
    },

    selectedStorage() {
      if (!this.item.source_storage_id || !this.secretStorages) return null;
      return this.secretStorages.find(s => s.id === this.item.source_storage_id);
    },
  },

  async created() {
    [
      this.secretStorages,
    ] = await Promise.all([
      this.loadProjectResources('secret_storages'),
    ]);

    // Load Azure configurations if user is admin
    if (this.user && this.user.admin) {
      try {
        const response = await this.$http.get('/api/azure');
        this.azureConfigs = response.data || [];
      } catch (error) {
        console.error('Failed to load Azure configurations:', error);
      }
    }
  },

  methods: {
    getNewItem() {
      return {
        ssh: {},
        login_password: {},
      };
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/keys`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/keys/${this.itemId}`;
    },

    async onAzureConfigChange(configId) {
      this.azureSubscriptions = [];
      this.azureKeyVaults = [];
      this.selectedSubscription = null;
      this.selectedKeyVault = null;

      if (!configId) return;

      this.loadingSubscriptions = true;
      try {
        const response = await this.$http.get(`/api/azure/${configId}/subscriptions`);
        this.azureSubscriptions = response.data || [];
      } catch (error) {
        console.error('Failed to load Azure subscriptions:', error);
        this.$emit('error', 'Failed to load Azure subscriptions');
      } finally {
        this.loadingSubscriptions = false;
      }
    },

    async onSubscriptionChange(subscriptionId) {
      this.azureKeyVaults = [];
      this.selectedKeyVault = null;

      if (!subscriptionId || !this.selectedAzureConfig) return;

      this.loadingKeyVaults = true;
      try {
        const response = await this.$http.get(`/api/azure/${this.selectedAzureConfig}/keyvaults?subscription_id=${subscriptionId}`);
        this.azureKeyVaults = response.data || [];
      } catch (error) {
        console.error('Failed to load Azure key vaults:', error);
        this.$emit('error', 'Failed to load Azure key vaults');
      } finally {
        this.loadingKeyVaults = false;
      }
    },

    onKeyVaultChange(vaultUrl) {
      // Update the source storage key with the selected vault URL
      if (vaultUrl) {
        this.item.source_storage_key = vaultUrl;
      }
    },
  },
};
</script>
