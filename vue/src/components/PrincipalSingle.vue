<template>
  <div class="rootbox">
    <h1>Principal Info</h1>
    <p>Name: {{ data.principal.name }}</p>

    <h2>Authorization</h2>
    <h3>Accounts</h3>
    <div class="box">
      <p class="generate-button">
        <b-button variant="primary" @click="onAddAccount">Add Account</b-button>
      </p>
      <template v-if="data.auths.length > 0">
        <data-table :items="data.auths" :fields="fields.auth" />
      </template>
      <template v-else>
        <p>No Authorizations.</p>
      </template>
    </div>

    <h3>API Keys</h3>
    <div class="box">
      <p class="generate-button">
        <b-button variant="primary" @click="onGenerateAPIKey">Generate API Key</b-button>
      </p>
      <template v-if="data.api_keys.length > 0">
        <data-table :items="data.api_keys" :fields="fields.api_key" />
      </template>
      <template v-else>
        <p>No API Keys.</p>
      </template>
    </div>

    <h2>Permission</h2>

    <h3>Roles</h3>
    <div class="box">
      <template v-if="data.roles.length > 0">
        <data-table :items="data.roles" :fields="fields.role" />
      </template>
      <template v-else>
        <p>No roles.</p>
      </template>
    </div>

    <h3>Groups</h3>
    <div class="box">
      <template v-if="data.groups.length > 0">
        <data-table :items="data.groups" :fields="fields.group" />
      </template>
      <template v-else>
        <p>No groups.</p>
      </template>
    </div>

  </div>
</template>

<script>
import axios from 'axios'
import DataTable from '@/components/DataTable'

export default {
  components: {
    'data-table': DataTable
  },

  data () {
    return {
      fields: {
        auth: ['account'],
        api_key: ['token'],
        role: ['name', 'description'],
        group: ['name', 'description']
      },
      data: {
        principal: {
          name: ''
        },
        auths: [],
        api_keys: [],
        roles: [],
        groups: []
      }
    }
  },

  created () {
    console.log('created')
    axios.get('/guardmech/api/principal/' + this.$route.params.id)
      .then(response => {
        console.log(response)
        this.data = response.data
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
  },

  methods: {
    onGenerateAPIKey: function () {
      console.log('Generate API Key!!')
    },
    onAddAccount: function () {
      console.log('Add Account')
    }
  }
}
</script>

<style scoped>
h1 {
    margin: 0.5em 0;
}

.box {
    margin: 20px 50px 0;
}

.generate-button {
    float: right;
    width: 200px;
    text-align: right;
}
</style>
