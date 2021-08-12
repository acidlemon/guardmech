<template>
  <div :class="'alert ' + alertVariant" role="alert">
    <slot />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'

type Props = {
   variant: string
}

export default defineComponent({
  name: 'BAlert',
  props: {
    variant: {
      type: String,
      default: '',
    },
  },
  setup(props: Props) {
    const alertVariant = computed(() => {
      if (!props.variant) {
        return ''
      }
      return `alert-${props.variant}`
    })

    return {
      alertVariant,
    }
  },
})
</script>
