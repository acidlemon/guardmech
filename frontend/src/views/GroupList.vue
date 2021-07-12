<template>
  <div class="container">
    <h2>Group List</h2>
    <section>
      <NewGroupModal />
    </section>
    <section>
      <BTable :data="groups" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/group/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import axios from 'axios'
import { Group } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import NewGroupModal from '@/components/modals/NewGroupModal.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewGroupModal,
  },
  setup() {
    const groups = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    onMounted(async () => {
      const res = await axios.get('/api/groups')
      groups.value = res.data.groups as Group[]
    })

    return {
      columns,
      groups,
    }
  },
})
</script>

