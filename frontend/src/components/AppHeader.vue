<template>
  <header class="header">
    <div class="container">
      <h2>Guardmech Control Panel</h2>
      <AppMenu />
    </div>
  </header>
</template>

<script>
import { defineComponent } from 'vue'
import AppMenu from '@/components/AppMenu.vue'

export default defineComponent({
  components: {
    AppMenu,
  },
  setup() {
    
  },
})
</script>

<style scoped>
.header {
  background: #F8F8F8;
  padding: 20px;

}
</style>
