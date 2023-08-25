<template>
  <div class="input">
    <label v-if="label" :for="targetId" class="form-label">{{ label }}</label>
    <div class="input-group">
      <span v-if="preText" class="input-group-text" id="basic-addon1">{{ preText }}</span>
      <input
        :type="type"
        class="form-control"
        :id="targetId"
        :placeholder="placeholder"
        :value="modelValue"
        @change="changed"
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext } from 'vue'

type InputType = 'text' | 'url' | 'tel' | 'search' | 'password' | 'number' | 'email'

type Props = {
  type: InputType
  label: string
  placeholder: string
  preText: string
}

export default defineComponent({
  name: 'BInput',
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    type: {
      type: String as () => InputType,
      default: 'text',
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    preText: {
      type: String,
      default: '',
    }
  },
  emits: ['update:modelValue'],
  setup(_: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)
    
    const changed = (e: Event) => {
      const el = e.currentTarget as HTMLInputElement
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
.input {
  padding-bottom: 10px;
}
</style>