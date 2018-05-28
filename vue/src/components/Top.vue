<template>
  <div class="hello">
    <h1>Guardmech Admin Panel</h1>
    <ul>
      <li><router-link to="/users">Users</router-link></li>
      <li><router-link to="/groups">Groups</router-link></li>
      <li><router-link to="/roles">Roles</router-link></li>
      <li><router-link to="/permissions">Permissions</router-link></li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
  text-align: left;
}
ul {
  text-align: left;
  margin: 20px;
  list-style-type: disc;
}
li {
  display: block;
}
a {
  color: #42b983;
}
</style>
