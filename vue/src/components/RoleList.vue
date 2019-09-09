<template>
  <div id="role_list">
    <h1>Roles</h1>

    <div class="rootbox">
      <data-table :items="roless" :fields="fields" path="/role/"/>
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
      fields: ['name', 'description', 'show_details'],
      roles: []
    }
  },

  created () {
    console.log('created')
    axios.get('/guardmech/api/roles')
      .then(response => {
        console.log(response)
        this.roles = response.data.roles
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
  }
}

</script>

<style>

h1 {
    margin: 1em;
}

.rootbox {
    margin: 20px 50px 0;
}

</style>
