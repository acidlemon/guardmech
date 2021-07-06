<template>
  <header>
    <h2>Guardmech Control Panel</h2>
    <AppMenu />
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
