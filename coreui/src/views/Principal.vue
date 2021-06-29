<template>
  <div class="animated fadeIn">
    <b-row>
      <b-col lg="12">
        <b-card>
          <div slot="header">
            Principal Info
          </div>
          <h2>Basic Info</h2>
          <b-table
            title="Basic Info"
            striped 
            :items="table_items.basic_info"
            :fields="table_fields.basic_info"
            class="detail-table"
          >
          </b-table>


          <h2>OpenID Connect Authorization Info</h2>
          <b-table
            striped
            :items="table_items.auths"
            :fields="table_fields.auths"
            class="detail-table"
          >
          </b-table>

          <p class="text-right floating-button">
            <b-button v-b-modal.new-token variant="danger">Create New</b-button>
            <b-modal
              id="new-token"
              title="Create API Token"
              hide-footer
              @show="resetAPIKeyModal"
            >
              <b-form  @submit="onNewAPIKey">
                <b-form-group id="input-group-1" label="API Key Name" label-for="input-1">
                  <b-form-input
                    id="input-1"
                    v-model="form.name"
                    required
                    placeholder="Enter name"
                  ></b-form-input>
                </b-form-group>
                <template v-if="form.token">
                  <b-form-group id="input-group-2" label="Issued Token" label-for="input-2">
                    <b-input-group>
                      <b-form-input
                        id="input-2"
                        type="text"
                        v-model="form.token"
                        disabled
                      ></b-form-input>
                      <b-input-group-append>
                        <b-button variant="info" @click="onCopyNewAPIKey">Copy</b-button>
                      </b-input-group-append>
                    </b-input-group>
                  </b-form-group>
                  <b-alert variant="warning" show>Note: You can copy this raw token now only. After closing modal, you cannot confirm it.</b-alert>
                  <b-button variant="secondary" @click="$bvModal.hide('new-token')">Close</b-button>
                </template>
                <template v-else-if="form.error">
                  <b-alert variant="danger" show>{{ form.error }}</b-alert>
                  <b-button variant="secondary" @click="$bvModal.hide('new-token')">Close</b-button>
                </template>
                <template v-else>
                  <b-button type="submit" variant="primary">Submit</b-button>
                </template>
              </b-form>
            </b-modal>
          </p>
          <h2>API Keys</h2>
          <b-table
            title="API Keys"
            striped
            :items="table_items.apikeys"
            :fields="table_fields.apikeys"
          >
            <template v-slot:cell(masked_token)="data">
              <span class="text-monospace">{{ data.value }}</span>
            </template>

          </b-table>

          <h2>Groups</h2>
          <b-table
            striped
            :items="table_items.groups"
            :fields="table_fields.groups"
            class="detail-table"
          >
          </b-table>

          <h2>Roles</h2>
          <b-table
            striped
            :items="table_items.roles"
            :fields="table_fields.roles"
            class="detail-table"
          >
          </b-table>

          <h2>Permissions</h2>
          <b-table
            striped
            :items="table_items.permissions"
            :fields="table_fields.permissions"
            class="detail-table"
          >
          </b-table>
        </b-card>
      </b-col>
    </b-row><!--/.row-->
  </div>    

</template>

<script>
import axios from 'axios'

async function fetchPrincipal(id) {
  const { data } = await axios.get('/guardmech/api/principal/' + id )

  return {
    basic_info: [
      { key: 'Name', value: data.principal.name},
      { key: 'Description', value: data.principal.description},
    ],
    auths: [data.auth_oidc],
    apikeys: data.auth_apikeys,
    groups: data.groups,
    roles: data.roles,
    permissions: data.permissions,
  }
}


export default {
  name: 'principal',
  props: [],
  data() {
    return {
      form: {
        name: "",
        token: "",
        error: "",
      },
      table_fields: {
        basic_info: [
          {key: 'key', label: 'Key'},
          {key: 'value', label: 'Value'},
        ],
        auths: [
          {key: 'issuer', label: 'OIDC Issuer'},
          {key: 'subject', label: 'OIDC Sub'},
          {key: 'email', label: 'Email'},
        ],
        apikeys: [
          {key: 'name', label: 'Token Name'},
          {key: 'masked_token', label: 'Token (Masked)'},
        ],
        groups: [
          {key: 'name', label: 'Group Name'},
          {key: 'description', label: 'Description'},
        ],
        roles: [
          {key: 'name', label: 'Role Name'},
          {key: 'description', label: 'Description'},
        ],
        permissions: [
          {key: 'name', label: 'Permission Name'},
          {key: 'description', label: 'Description'},
        ],
              },
      table_items: {
        basic_info: [],
        auths: [],
        apikeys: [],
        groups: [],
        roles: [],
        permissions: [],
      },
    }
  },
  async mounted() {
    this.table_items = await fetchPrincipal(this.$route.params.id)
  },
  methods: {
    onCopyNewAPIKey() {
      navigator.clipboard.writeText(this.form.token)
    },
    resetAPIKeyModal() {
      this.form = {
        name: "",
        token: "",
        error: "",
      }
    },
    async onNewAPIKey(evt) {
      evt.preventDefault()
      let params = new URLSearchParams()
      params.append('name', this.form.name)
      axios.post('/guardmech/api/principal/' + this.$route.params.id + '/new_key', params).then(async (response) => {
        console.log(response)
        this.table_items = await fetchPrincipal(this.$route.params.id)
        this.form.token = response.data.token

      }).catch(error => {
        this.form.error = error
        console.log(error)
      })
    }
  }
}
</script>

<style scoped>
.floating-button {
  float: right;
}
.detail-table {
  margin-bottom: 40px;
}
</style>

undefined