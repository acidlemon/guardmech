<template>
  <div class="select">
    <label v-if="label" :for="targetId" class="form-label">{{ label }}</label>
    <select class="form-select" :id="targetId" @change="changed">
      <option
        v-for="item in items"
        :value="item.value"
        :key="item.label"
        :disabled="item.disabled"
      >{{ item.label }}</option>
    </select>
    <div v-if="selectedTips" class="tips">{{ selectedTips }}</div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, SetupContext } from 'vue'

export type BSelectItem = {
  label: string
  value: string
  tips?: string
  disabled?: boolean
}

type Props = {
  modelValue: string
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
  setup(props: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)

    const changed = (e: Event) => {
      const el = e.currentTarget as HTMLSelectElement
      context.emit('update:modelValue', el.value)
    }

    const selectedTips = computed<string>(() => {
      const item = props.items.find(x => x.value === props.modelValue)
      if (!item || !item.tips) return ''

      return item.tips
    })

    return {
      targetId,
      changed,
      selectedTips,
    }
  },
})
</script>

<style scoped>
.select {
  padding-bottom: 10px;
}
.tips {
  padding: 5px 0 5px 5px;
  font-size: 0.8em;
  color: gray;
}
</style>