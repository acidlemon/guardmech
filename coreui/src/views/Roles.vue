<template>
  <div class="animated fadeIn">
    <b-row>
      <b-col lg="12">
        <b-card>
          <div slot="header">
            Roles
          </div>
          <b-table :items="items" :fields="fields">
            <template v-slot:cell(editable)="data">
              <b-button variant="primary">編集</b-button>
            </template>
          </b-table>
        </b-card>
      </b-col>
    </b-row><!--/.row-->
  </div>    
</template>

<script>
import axios from 'axios'

export default {
  name: 'roles',
  data() {
    return {
      items: [
      ],
      fields: [
        {key: 'id'},
        {key: 'name'},
        {key: 'description'},
        {key: 'editable', label: 'hoge'},
      ]
    }
  },
  mounted() {
    axios.get('/guardmech/api/roles').then(response => {
      let data = response.data.roles
      //for (let d of data) {
      //  console.log(d)
      //  d["editable"] = true
      //}
      this.items = data
    })

  }
}
</script>

