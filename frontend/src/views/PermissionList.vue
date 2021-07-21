<template>
  <div class="container">
    <h2>Permission List</h2>
    <section>
      <NewPermissionModal @completed="created"/>
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
import NewPermissionModal from '@/components/modals/NewPermissionModal.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewPermissionModal,
  },
  setup() {
    const permissions = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    const fetchList = (async () => {
      const res = await axios.get('/api/permissions').catch(e => e.response)
      permissions.value = res.data.permissions as Permission[]
    })

    onMounted(async () => {
      fetchList()
    })

    const created = (() => {
      fetchList()
    })

    return {
      columns,
      permissions,
      created,
    }
  },
})
</script>

