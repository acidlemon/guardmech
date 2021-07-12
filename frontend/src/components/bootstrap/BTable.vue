<template>
  <table class="table table-striped">
    <thead>
      <tr>
        <th scope="col" v-for="h in columns" :key="h.key">{{ h.label }}</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="row in data" :key="row">
        <td v-for="h in columns" :key="h.key">
          <template v-if="slotAccessor(h.key)">
            <slot :name="`cell(${h.key})`" :row="row" />
          </template>
          <template v-else>
            {{ row[h.key] }}
          </template>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { computed, defineComponent, SetupContext } from 'vue'

export type BTableColumn = {
  key: string
  label: string
}

export type BTableRow = {
  [index: string]: any
}

type Props = {
  data: BTableRow[]
  columns: BTableColumn[]
  variant: string
}

export default defineComponent({
  name: 'BTable',
  props: {
    data: {
      type: Array as () => BTableRow[],
      default: () => [],
    },
    columns: {
      type: Array as () => BTableColumn[],
      default: () => [],
    },
    variant: {
      type: String,
      default: '',
    }
  },
  setup(props: Props, context: SetupContext) {
    const tableVariant = computed(() => props.variant ? `table-${props.variant}` : 'table')

    const slot = computed(() => context.slots)

    const slotAccessor = (key: string) => {
      const slotData = slot.value[`cell(${key})`]
      if (!slotData) { return undefined }
      return slotData()
    }

    return {
      slot,
      slotAccessor,
      tableVariant
    }
  },
})
</script>
