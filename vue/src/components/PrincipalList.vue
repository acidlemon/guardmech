<template>
  <div id="principal_list">
    <h1>Principals</h1>

    <div class="rootbox">
      <data-table :items="principals" :fields="fields" path="/principal/"/>
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
      principals: []
    }
  },

  created () {
    console.log('created')
    axios.get('/guardmech/api/principals')
      .then(response => {
        console.log(response)
        this.principals = response.data.principals
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
