<template>
  <div class="container">
    <h2>Group List</h2>
    <AuthorityStatusBox :status="authorityStatus" />
    <section v-if="canWrite">
      <NewGroupModal @completed="created" />
    </section>
    <section v-if="canRead">
      <BTable :data="groups" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/group/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import axios from 'axios'
import { Group } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import NewGroupModal from '@/components/modals/NewGroupModal.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow, BTableColumn } from '@/types/bootstrap'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewGroupModal,
    AuthorityStatusBox,
  },
  setup() {
    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()
    
    const groups = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    const fetchList = (async () => {
      const res = await axios.get('/api/groups').catch(e => e.response)
      groups.value = res.data.groups as Group[]
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
      groups,
      created,
    }
  },
})
</script>

