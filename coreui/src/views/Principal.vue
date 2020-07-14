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


          <b-card>
            <template v-slot:header>
              <h4 class="mb-0">Authorizations</h4>
            </template>
            <b-table
              striped
              :items="table_items.auths"
              :fields="table_fields.auths"
              class="detail-table"
            >
            </b-table>
          </b-card>

          <p class="text-right floating-button">
            <b-button  v-b-modal.new-token variant="danger">Create New</b-button>
            <b-modal id="new-token" title="Create API Token" hide-footer>
              <b-form>
                <b-form-group id="input-group-1" label="API Key Name" label-for="input-1">
                  <b-form-input
                    id="input-1"
                    v-model="form.name"
                    required
                    placeholder="Enter name"
                  ></b-form-input>
                </b-form-group>
                <b-form-group id="input-group-1" label="Description">
                  <b-form-input
                    id="input-2"
                    v-model="form.description"
                  ></b-form-input>
                  </b-form-group>
                <b-button variant="primary" @click="onNewAPIKey">Submit</b-button>
              </b-form>
            </b-modal>
          </p>
          <h2>API Keys</h2>
          <b-table title="API Keys" striped :items="table_items.api_keys" :fields="table_fields.api_keys">
          </b-table>
        </b-card>
      </b-col>
    </b-row><!--/.row-->
  </div>    

</template>

<script>
import axios from 'axios'

export default {
  name: 'principal',
  props: [],
  data() {
    return {
      form: {
        name: "",
        description: "",
      },
      table_fields: {
        basic_info: [
          {key: 'key', label: 'Key'},
          {key: 'value', label: 'Value'},
        ],
        auths: [
          {key: 'unique_id', label: 'Unique ID'},
          {key: 'issuer', label: 'OIDC Issuer'},
          {key: 'subject', label: 'OIDC Sub'},
          {key: 'email', label: 'Email'},
        ],
        api_keys: [
          {key: 'unique_id', label: 'Unique ID'},
          {key: 'name', label: 'Token Name'},
          {key: 'masked_token', label: 'Token (Masked)'},
        ],
      },
      table_items: {
        basic_info: [],
        auths: [],
        groups: [],
        basic_info: [],
      },
    }
  },
  mounted() {
    this.fetchPrincipal()
  },
  methods: {
    fetchPrincipal() {
      axios.get('/guardmech/api/principal/' + this.$route.params.seq_id ).then(response => {
        console.log(response)

        // basic info
        const pri = response.data.principal
        this.table_items.basic_info.push({ key: 'Unique ID', value: pri.unique_id})
        this.table_items.basic_info.push({ key: 'Name', value: pri.name})
        this.table_items.basic_info.push({ key: 'Description', value: pri.description})

        // auths
        this.table_items.auths = response.data.auths
        this.table_items.api_keys = response.data.api_keys
      })
    },
    onNewAPIKey() {
      console.log('onNewAPIKey')
      let params = new URLSearchParams()
      params.append('name', this.form.name)
      params.append('description', this.form.description)
      axios.post('/guardmech/api/principal', params).then(response => {
        console.log(response)
        fetchPrincipal()
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