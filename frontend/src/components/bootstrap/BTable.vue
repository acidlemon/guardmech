<script setup lang="ts">
import type { BTableColumn, BTableRow } from '@/types/bootstrap'
import { useSlots } from 'vue'
type Props = {
  data: BTableRow[]
  columns: BTableColumn[]
  variant?: string
}
const _ = withDefaults(defineProps<Props>(), {
  data: () => [],
  columns: () => [],
})
const slots = useSlots()

//const tableVariant = computed(() => props.variant ? `table-${props.variant}` : 'table')
const slotAccessor = (key: string) => {
  const slotData = slots[`cell(${key})`]
  if (!slotData) { return }
  return slotData()
}
</script>

<template>
  <table class="table table-striped">
    <thead>
      <tr>
        <th
          v-for="h in columns"
          :key="h.key"
          scope="col"
        >
          {{ h.label }}
        </th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="row in data"
        :key="row['id']"
      >
        <td
          v-for="(h, index) in columns"
          :key="h.key"
        >
          <template v-if="slotAccessor(h.key)">
            <slot
              :name="`cell(${h.key})`"
              :index="index"
              :row="row"
            />
          </template>
          <template v-else>
            {{ row[h.key] }}
          </template>
        </td>
      </tr>
    </tbody>
  </table>
</template>