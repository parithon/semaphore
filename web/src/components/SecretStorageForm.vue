<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="item != null"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-text-field
        v-model="item.name"
        :label="$t('name')"
        :rules="[v => !!v || $t('name_required')]"
        required
        :disabled="formSaving"
        outlined
        dense
    ></v-text-field>

    <v-select
        v-model="item.type"
        :label="$t('type')"
        :rules="[v => !!v || $t('type_required')]"
        :items="secretStorageTypes"
        item-value="id"
        item-text="name"
        required
        :disabled="formSaving"
        outlined
        dense
    />

    <div v-if="item.type === 'vault'">

      <v-text-field
          v-model="item.params.url"
          :label="$t('Server URL')"
          :disabled="formSaving"
          :rules="[v => !!v || $t('url_required')]"
          required
          data-testid="secretStorage-vaultURL"
          outlined
          dense
      ></v-text-field>

      <v-text-field
          v-model="item.vault_token"
          :label="$t('Token')"
          :disabled="formSaving"
          :rules="[v => !!v || itemId !== 'new' || $t('token_required')]"
          required
          data-testid="secretStorage-vaultToken"
          outlined
          dense
          append-icon="mdi-lock"
      ></v-text-field>

      <v-checkbox
        v-model="item.readonly"
        :label="$t('Read only')"
        :disabled="formSaving"
      />
    </div>

    <div v-if="item.type === 'azure'">

      <v-text-field
          v-model="item.params.vault_url"
          :label="$t('Vault URL')"
          :disabled="formSaving"
          :rules="[v => !!v || $t('vault_url_required')]"
          required
          data-testid="secretStorage-azureVaultURL"
          outlined
          dense
          hint="e.g., https://your-keyvault.vault.azure.net/"
      ></v-text-field>

      <v-text-field
          v-model="item.params.tenant_id"
          :label="$t('Tenant ID')"
          :disabled="formSaving"
          :rules="[v => !!v || $t('tenant_id_required')]"
          required
          data-testid="secretStorage-azureTenantID"
          outlined
          dense
      ></v-text-field>

      <v-text-field
          v-model="item.params.client_id"
          :label="$t('Client ID')"
          :disabled="formSaving"
          :rules="[v => !!v || $t('client_id_required')]"
          required
          data-testid="secretStorage-azureClientID"
          outlined
          dense
      ></v-text-field>

      <v-text-field
          v-model="item.params.client_secret"
          :label="$t('Client Secret')"
          :disabled="formSaving"
          :rules="[v => !!v || $t('client_secret_required')]"
          required
          data-testid="secretStorage-azureClientSecret"
          outlined
          dense
          type="password"
          append-icon="mdi-lock"
      ></v-text-field>

      <v-checkbox
        v-model="item.readonly"
        :label="$t('Read only')"
        :disabled="formSaving"
      />
    </div>
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  props: {
    itemType: String,
  },

  mixins: [ItemFormBase],

  data() {
    return {
      secretStorageTypes: [{
        id: 'vault',
        name: 'Hashicorp Vault',
      }, {
        id: 'azure',
        name: 'Azure Key Vault',
      }],
    };
  },

  methods: {

    getNewItem() {
      return {
        params: {},
      };
    },

    afterLoadData() {
      if (!this.item.params) {
        this.item.params = {};
      }

      this.item.type = this.itemType;
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/secret_storages`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/secret_storages/${this.itemId}`;
    },
  },
};
</script>
