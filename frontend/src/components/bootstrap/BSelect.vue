<template>
  <div class="select">
    <label v-if="label" :for="targetId" class="form-label">{{ label }}</label>
    <select class="form-select" :id="targetId" @change="changed">
      <option v-for="item in items" :value="item.value" :key="item.label">{{ item.label }}</option>
    </select>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext } from 'vue'

export type BSelectItem = {
  label: string
  value: string
}

type Props = {
  label: string
  items: BSelectItem[]
}

export default defineComponent({
  name: 'BSelect',
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    items: {
      type: Array as () => BSelectItem[],
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
  },
  emits: ['update:modelValue'],
  setup(_: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)

    const changed = (e: Event) => {
      const el = e.currentTarget as HTMLSelectElement
      context.emit('update:modelValue', el.value)
    }

    return {
      targetId,
      changed,
    }
  },
})
</script>

<style scoped>
.select {
  padding-bottom: 10px;
}
</style>