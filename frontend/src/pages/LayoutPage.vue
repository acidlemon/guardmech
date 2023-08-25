<template>
  <div class="app">
    <AppHeader />
    <router-view />
  </div>
</template>

<script lang="ts">
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
