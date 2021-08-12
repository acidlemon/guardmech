<template>
  <div class="container">
    <h2>Principal List</h2>
    <AuthorityStatusBox :status="authorityStatus" />
    <section v-if="canWrite">
      <NewPrincipalModal @completed="created"/>
    </section>
    <section v-if="canRead">
      <BTable :data="principals" :columns="columns">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/principal/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import axios from 'axios'
import { PrincipalPayload } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import NewPrincipalModal from '@/components/modals/NewPrincipalModal.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewPrincipalModal,
    AuthorityStatusBox,
  },
  setup() {
    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

    const principals = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    const fetchList = (async () => {
      const res = await axios.get('/api/principals').catch(e => e.response)
      const payload = res.data.principals as PrincipalPayload[]
      principals.value = payload.map(d => d.principal)
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchList()
        }
      }
    })

    const created = (() => {
      fetchList()
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      columns,
      principals,
      created,
    }
  },
})
</script>
