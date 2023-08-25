<template>
  <div class="check-list">
    <label v-if="label" class="form-label">{{ label }}</label>
    <div class="form-check" v-for="(item, index) in items" :key="item.value">
      <input
        class="form-check-input"
        :type="checkType"
        :value="item.value"
        :id="`${targetId}-${index}`"
        :name="name"
        @change="changed"
        :checked="item.defaultChecked"
      >
      <label class="form-check-label" :for="`${targetId}-${index}`">
        {{ item.label }}
      </label>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext } from 'vue'

type CheckBoxType = 'checkbox' | 'radio'

export type BCheckItem = {
  label: string
  value: string
  defaultChecked: boolean
}

type Props = {
  modelValue: string
  checkType: CheckBoxType
  items: BCheckItem[]
  name: string
  label: string
}

export default defineComponent({
  name: 'BCheckList',
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    checkType: {
      type: String as () => CheckBoxType,
      default: 'checkbox',
    },
    items: {
      type: Array as () => BCheckItem[],
      default: () => [],
    },
    name: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
  },
  setup(_: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)

    const changed = (e: Event) => {
      const el = e.currentTarget as HTMLInputElement
      context.emit('update:modelValue', el.value)
    }
    
    return {
      targetId,
      changed
    }
  },
})
</script>

<style scoped>
.check-list {
  padding-bottom: 10px;
}
.form-check {
  margin-left: 20px;
}
</style>