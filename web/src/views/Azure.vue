<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items != null">
    <EditDialog
      v-model="editDialog"
      save-button-text="Save"
      title="Azure Configuration"
      @save="loadItems()"
    >
      <template v-slot:form="{ onSave, onError, needSave, needReset }">
        <AzureForm
          :item-id="itemId"
          @save="onSave"
          @error="onError"
          :need-save="needSave"
          :need-reset="needReset"
        />
      </template>
    </EditDialog>

    <YesNoDialog
      title="Delete Azure Configuration"
      text="Are you sure you want to delete this Azure configuration?"
      v-model="deleteItemDialog"
      @yes="deleteItem(itemId)"
    />

    <v-toolbar flat>
      <v-btn
        icon
        class="mr-4"
        @click="returnToProjects()"
      >
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>Azure Configurations</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        @click="editItem('new')"
      >New Azure Configuration</v-btn>
    </v-toolbar>

    <v-divider />

    <v-data-table
      :headers="headers"
      :items="items"
      class="mt-4"
      :footer-props="{ itemsPerPageOptions: [20] }"
    >
      <template v-slot:item.created="{ item }">
        {{ formatDate(item.created) }}
      </template>

      <template v-slot:item.actions="{ item }">
        <div style="white-space: nowrap">
          <v-btn
            icon
            class="mr-1"
            @click="askDeleteItem(item.id)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            icon
            class="mr-1"
            @click="editItem(item.id)"
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
        </div>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import EventBus from '@/event-bus';
import YesNoDialog from '@/components/YesNoDialog.vue';
import ItemListPageBase from '@/components/ItemListPageBase';
import EditDialog from '@/components/EditDialog.vue';
import AzureForm from '@/components/AzureForm.vue';

export default {
  mixins: [ItemListPageBase],

  components: {
    YesNoDialog,
    AzureForm,
    EditDialog,
  },

  methods: {
    getHeaders() {
      return [{
        text: 'Name',
        value: 'name',
        width: '30%',
      },
      {
        text: 'Description',
        value: 'description',
        width: '40%',
      },
      {
        text: 'Tenant ID',
        value: 'tenant_id',
        width: '20%',
      },
      {
        text: 'Created',
        value: 'created',
        width: '15%',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
        width: '10%',
      }];
    },

    async returnToProjects() {
      EventBus.$emit('i-open-last-project');
    },

    getItemsUrl() {
      return '/api/azure';
    },

    getSingleItemUrl() {
      return `/api/azure/${this.itemId}`;
    },

    getEventName() {
      return 'i-azure';
    },

    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
    },
  },
};
</script>