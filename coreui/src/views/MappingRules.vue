<template>
  <div class="animated fadeIn">
    <b-row>
      <b-col lg="12">
        <b-card>
          <div slot="header">
            Mapping Rules
          </div>
          <p class="text-right">
            <b-button  v-b-modal.new-mapping-rule variant="danger">Create New Rule</b-button>
            <b-modal id="new-mapping-rule" title="Add New Mapping Rule" hide-footer>
              <b-form> <!--@submit="onSubmit" @reset="onReset"-->
                <b-form-group id="input-group-1" label="Rule Name:" label-for="input-1">
                  <b-form-input
                    id="input-1"
                    v-model="form.name"
                    required
                    placeholder="Name"
                  ></b-form-input>
                </b-form-group>

                <b-form-group id="input-group-2" label="Description:" label-for="input-2">
                  <b-form-input
                    id="input-2"
                    v-model="form.description"
                    description="Write your memo"
                    placeholder="My Group"
                  ></b-form-input>
                </b-form-group>

                <b-form-group id="input-group-3">
                  <b-form-radio-group v-model="form.checked" id="checkboxes-3">
                    <b-form-radio value="domain">Mail Address Host</b-form-radio>
                    <b-form-radio value="group">Provided Group</b-form-radio>
                  </b-form-radio-group>
                </b-form-group>


                <b-form-group
                  id="input-group-4"
                  label="Email address:"
                  label-for="input-1"
                >
                  <b-form-input
                    id="input-4"
                    v-model="form.param"
                    type="email"
                    required
                    placeholder="@example.com"
                  ></b-form-input>
                </b-form-group>

                <b-button type="submit" variant="primary">Submit</b-button>
              </b-form>
            </b-modal>
          </p>
          <b-table :items="items" :fields="fields">
            <template v-slot:cell(editable)="">
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
  name: 'mapping-rules',
  data() {
    return {
      items: [
      ],
      fields: [
        {key: 'id'},
        {key: 'name'},
        {key: 'description'},
        {key: 'editable', label: 'hoge'},
      ],
      form: {
        name: "",
        description: "",
        email: "",
      },
    }
  },
  mounted() {
    axios.get('/guardmech/api/mapping_rules').then(response => {
      let data = response.data.permissions
      //for (let d of data) {
      //  console.log(d)
      //  d["editable"] = true
      //}
      this.items = data
    })

  }
}
</script>

