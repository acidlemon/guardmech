<template>
  <div>
    <h2>Permission List</h2>
    <section>
      <BButton variant="danger">Create New</BButton>
    </section>
    <section>
      <BTable :data="permissions" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/permission/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import axios from 'axios'
import { Permission } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
  },
  setup() {
    const permissions = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    onMounted(async () => {
      const res = await axios.get('/api/permissions')
      permissions.value = res.data.permissions as Permission[]
    })

    return {
      columns,
      permissions,
    }
  },
})
</script>

