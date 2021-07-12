<template>
  <div class="container">
    <h2>Principal List</h2>
    <section>
      <NewPrincipalModal @completed="created"/>
    </section>
    <section>
      <BTable :data="principals" :columns="columns">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/principal/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, defineComponent } from 'vue'
import axios from 'axios'
import { PrincipalPayload } from '@/types/api'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import NewPrincipalModal from '@/components/modals/NewPrincipalModal.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewPrincipalModal,
  },
  setup() {
    const principals = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    const fetchList = (async () => {
      const res = await axios.get('/api/principals')
      const payload = res.data.principals as PrincipalPayload[]
      principals.value = payload.map(d => d.principal)
    })

    onMounted(() => {
      fetchList()
    })

    const created = (() => {
      fetchList()
    })

    return {
      columns,
      principals,
      created,
    }
  },
})
</script>
