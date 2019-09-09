<template>
  <b-table striped hover :items="items" :fields="fields">
    <template slot="show_details" slot-scope="data">
      <b-button :to=" path + data.item.id ">Show</b-button>
    </template>
  </b-table>
</template>

<script>
export default {
  props: [ 'items', 'fields', 'path' ]
}
</script>

<style>

</style>
