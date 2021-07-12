<template>
  <div class="container">
    <h2>Mapping Rule List</h2>
    <section>
      <BButton variant="danger">Create New</BButton>
    </section>
    <section>
      <BTable
        v-if="mappingRules.length"
        :data="mappingRules"
        :columns="columns"
        variant="primary"
      >
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/mapping_rule/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
      <p>There's no mapping rules.</p>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import axios from 'axios'
import { MappingRule } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
  },
  setup() {
    const mappingRules = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    onMounted(async () => {
      const res = await axios.get('/api/mapping_rules')
      mappingRules.value = res.data.mapping_rules as MappingRule[]
    })

    return {
      columns,
      mappingRules,
    }
  },
})
</script>
