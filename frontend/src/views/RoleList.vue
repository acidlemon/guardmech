<template>
  <div class="container">
    <h2>Role List</h2>
    <AuthorityStatusBox :status="authorityStatus" />
    <section v-if="canWrite">
      <NewRoleModal @completed="created"/>
    </section>
    <section v-if="canRead">
      <BTable :data="roles" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/role/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import axios from 'axios'
import { Role } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import NewRoleModal from '@/components/modals/NewRoleModal.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow, BTableColumn } from '@/types/bootstrap'


export default defineComponent({
  components: {
    BButton,
    BTable,
    NewRoleModal,
    AuthorityStatusBox,
  },
  setup() {
    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

    const roles = ref<BTableRow[]>([])
    const columns = ref<BTableColumn[]>([
      { key: 'name', label: 'Name' },
      { key: 'description', label: 'Description' },
      { key: 'action', label: '' },
    ])

    const fetchList = (async () => {
      const res = await axios.get('/api/roles').catch(e => e.response)
      roles.value = res.data.roles as Role[]
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
      roles,
      created,
    }
  },
})
</script>

