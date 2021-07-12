<template>
  <div class="container">
    <h2>Role List</h2>
    <section>
      <NewRoleModal />
    </section>
    <section>
      <BTable :data="roles" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/role/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import axios from 'axios'
import { Role } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import NewRoleModal from '@/components/modals/NewRoleModal.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewRoleModal,
  },
  setup() {
    const roles = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    onMounted(async () => {
      const res = await axios.get('/api/roles')
      roles.value = res.data.roles as Role[]
    })

    return {
      columns,
      roles,
    }
  },
})
</script>

