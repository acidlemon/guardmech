<template>
  <span>
    <a
      v-if="href"
      :href="href"
      :class="`btn ${variantClass}`"
      role="button"
    >
      <slot />
    </a>

    <router-link
      v-else-if="link"
      :to="link"
    >
      <button
        type="button"
        :class="`btn ${variantClass}`"
      >
        <slot />
      </button>
    </router-link>

    <button
      v-else
      type="button"
      :data-bs-toggle="toggle"
      :class="`btn ${variantClass}`"
      @click="clicked"
    >
      <slot />
    </button>
  </span>
</template>

<script lang="ts">
import { computed, defineComponent, SetupContext } from 'vue'

type Props = {
  href: string
  link: string
  variant: string
  toggle: string
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
    toggle: {
      type: String,
      default: '',
    },
  },
  emits: ['click'],
  setup(props: Props, context: SetupContext) {
    const variantClass = computed(() => {
      if (!props.variant) { return 'btn-primary'}
      return `btn-${props.variant}`
    })

    const clicked = (() => {
      context.emit('click')
    })

    return {
      variantClass,
      clicked,
    }
  },
})
</script>
