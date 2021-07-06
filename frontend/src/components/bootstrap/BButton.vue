<template>
  <span>
    <a v-if="href" :href="href" :class="`btn ${variantClass}`" role="button"><slot /></a>
    <router-link v-else-if="link" :to="link"><button type="button" :class="`btn ${variantClass}`"><slot /></button></router-link>
    <button type="button" :class="`btn ${variantClass}`" v-else><slot /></button>
  </span>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'

type Props = {
  href: string
  link: string
  variant: string
}

export default defineComponent({
  name: 'BButton',
  props: {
    href: {
      type: String,
      default: '',
    },
    link: {
      type: String,
      default: '',
    },
    variant: {
      type: String,
      default: '',
    },
  },
  setup(props: Props) {
    const variantClass = computed(() => {
      if (!props.variant) { return 'btn-primary'}
      return `btn-${props.variant}`
    })
    return {variantClass}
  },
})
</script>
