<template>
  <div class="app">
    <AppHeader />
    <router-view />
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import AppHeader from '@/components/AppHeader.vue'

export default defineComponent({
  components: {
    AppHeader,
  },
  setup() {
  },
})
</script>
