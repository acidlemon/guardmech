<template>
  <header class="header">
    <div class="container">
      <h2>Guardmech Control Panel</h2>
      <AppMenu />
    </div>
  </header>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import AppMenu from '@/components/AppMenu.vue'

export default defineComponent({
  components: {
    AppMenu,
  },
  setup() {
    return {}
  },
})
</script>

<style scoped>
.header {
  background: #F8F8F8;
  padding: 20px;
  margin-bottom: 20px;
}
</style>
