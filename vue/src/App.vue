
<template>
<div id="app">
  <app-header />
  <div id="main">
    <router-view/>
  </div>
  <app-footer />
</div>
</template>

<script>
import Header from '@/components/Header'
import Footer from '@/components/Footer'

export default {
  name: 'App',
  components: {
    'app-header': Header,
    'app-footer': Footer
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  text-align: center;
}

#main {
  padding: 20px 50px;
  text-align: left;
}

h1 {
    margin: 0;
    padding: 0;
}
</style>
