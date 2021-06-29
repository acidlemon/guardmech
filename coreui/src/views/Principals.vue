<template>
  <div class="animated fadeIn">
    <b-row>
      <b-col lg="12">
        <b-card>
          <div slot="header">
            Principals
          </div>
          <p class="text-right">
            <b-button  v-b-modal.new-principal variant="danger">Create New</b-button>
            <b-modal id="new-principal" title="Add New Principal" hide-footer>
              <b-form> <!--@submit="onSubmit" @reset="onReset"-->
                <b-form-group id="input-group-1" label="Your Name:" label-for="input-1">
                  <b-form-input
                    id="input-2"
                    v-model="form.name"
                    required
                    placeholder="Enter name"
                  ></b-form-input>
                </b-form-group>

                <b-form-group id="input-group-2" label="Description:" label-for="input-2">
                  <b-form-input
                    id="input-2"
                    v-model="form.description"
                    required
                    placeholder="Write your memo"
                  ></b-form-input>
                </b-form-group>

                <b-form-group id="input-group-3">
                  <b-form-checkbox-group v-model="form.checked" id="checkboxes-3">
                    <b-form-checkbox value="me">Check me out</b-form-checkbox>
                  </b-form-checkbox-group>
                </b-form-group>


                <b-form-group
                  id="input-group-4"
                  label="Email address:"
                  label-for="input-1"
                  description="We'll never share your email with anyone else."
                >
                  <b-form-input
                    id="input-4"
                    v-model="form.email"
                    type="email"
                    required
                    placeholder="Enter email"
                  ></b-form-input>
                </b-form-group>

                <b-button type="submit" variant="primary">Submit</b-button>
              </b-form>
            </b-modal>
          </p>
          <b-table :items="items" :fields="fields">
            <template v-slot:cell(action)="data">
              <b-button variant="primary" :to="'principal/' + data.item.id">情報</b-button>
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
  name: 'principals',
  data() {
    return {
      items: [
      ],
      fields: [
        {key: 'name'},
        {key: 'description'},
        {key: 'action', label: ''},
      ],
      form: {
        email: "",
        name: "",
      },
      foods: ["a", "b"],
    }
  },
  mounted() {
    axios.get('/guardmech/api/principals').then(response => {
      let data = response.data.principals
      //for (let d of data) {
      //  console.log(d)
      //  d["editable"] = true
      //}
      this.items = data.map(d => d.principal)
    })

  }
}
</script>

