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
      label="Name"
      :rules="[v => !!v || 'Name is required']"
      required
      :disabled="formSaving"
      outlined
      dense
    ></v-text-field>

    <v-textarea
      v-model="item.description"
      label="Description (Optional)"
      :disabled="formSaving"
      outlined
      dense
      rows="2"
    ></v-textarea>

    <v-text-field
      v-model="item.tenant_id"
      label="Tenant ID"
      :rules="[v => !!v || 'Tenant ID is required']"
      required
      :disabled="formSaving"
      outlined
      dense
      placeholder="e.g., 12345678-1234-1234-1234-123456789012"
    ></v-text-field>

    <v-text-field
      v-model="item.client_id"
      label="Client ID"
      :rules="[v => !!v || 'Client ID is required']"
      required
      :disabled="formSaving"
      outlined
      dense
      placeholder="e.g., 87654321-4321-4321-4321-210987654321"
    ></v-text-field>

    <v-text-field
      v-model="item.client_secret"
      :append-icon="showClientSecret ? 'mdi-eye' : 'mdi-eye-off'"
      label="Client Secret"
      :rules="[v => !!v || 'Client Secret is required']"
      :type="showClientSecret ? 'text' : 'password'"
      required
      :disabled="formSaving"
      autocomplete="new-password"
      @click:append="showClientSecret = !showClientSecret"
      outlined
      dense
      placeholder="Enter your Azure service principal client secret"
    ></v-text-field>

    <v-alert
      dense
      text
      type="info"
      class="mt-4"
    >
      <div class="subtitle-2 mb-2">Azure Setup Requirements:</div>
      <ul class="pl-4">
        <li>Create a service principal in Azure Active Directory</li>
        <li>Grant appropriate permissions to access Key Vault and subscriptions</li>
        <li>Required permissions: Key Vault Secrets Officer or equivalent</li>
        <li>Ensure the service principal has access to the subscriptions you want to use</li>
      </ul>
    </v-alert>

    <div v-if="!isNew && item.id" class="mt-4">
      <v-expansion-panels>
        <v-expansion-panel>
          <v-expansion-panel-header>
            <div>
              <v-icon left>mdi-microsoft-azure</v-icon>
              Test Azure Connection
            </div>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
            <v-btn
              color="primary"
              :loading="testingConnection"
              @click="testConnection"
              class="mb-3"
            >
              Test Connection
            </v-btn>
            
            <div v-if="subscriptions.length > 0">
              <div class="subtitle-2 mb-2">Available Subscriptions:</div>
              <v-list dense>
                <v-list-item v-for="sub in subscriptions" :key="sub.id">
                  <v-list-item-content>
                    <v-list-item-title>{{ sub.name }}</v-list-item-title>
                    <v-list-item-subtitle>{{ sub.id }} ({{ sub.state }})</v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </div>
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
  </v-form>
</template>

<script>
import ItemFormBase from '@/components/ItemFormBase';
import axios from 'axios';

export default {
  mixins: [ItemFormBase],

  data() {
    return {
      showClientSecret: false,
      testingConnection: false,
      subscriptions: [],
    };
  },

  methods: {
    getNewItem() {
      return {
        name: '',
        description: '',
        tenant_id: '',
        client_id: '',
        client_secret: '',
      };
    },

    getItemsUrl() {
      return '/api/azure';
    },

    getSingleItemUrl() {
      return `/api/azure/${this.itemId}`;
    },

    async testConnection() {
      if (!this.item.id) {
        this.$emit('error', 'Please save the configuration first');
        return;
      }

      this.testingConnection = true;
      try {
        const response = await axios.get(`/api/azure/${this.item.id}/subscriptions`);
        this.subscriptions = response.data || [];
        
        if (this.subscriptions.length === 0) {
          this.$emit('error', 'No subscriptions found. Please check your Azure configuration.');
        } else {
          this.$emit('success', `Successfully connected! Found ${this.subscriptions.length} subscription(s).`);
        }
      } catch (error) {
        console.error('Azure connection test failed:', error);
        this.$emit('error', error.response?.data?.message || 'Failed to connect to Azure');
      } finally {
        this.testingConnection = false;
      }
    },
  },
};
</script>