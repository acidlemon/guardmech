<template>
  <div class="container">
    <h2>Permission List</h2>
    <AuthorityStatusBox :status="authorityStatus" />
    <section v-if="canWrite">
      <NewPermissionModal @completed="created"/>
    </section>
    <section v-if="canRead">
      <BTable :data="permissions" :columns="columns" variant="primary">
        <template #cell(action)="data">
          <BButton v-if="data?.row" :link="`/permission/${data.row.id}`" >View</BButton>
        </template>
      </BTable>
    </section>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import axios from 'axios'
import { Permission } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import BButton from '@/components/bootstrap/BButton.vue'
import NewPermissionModal from '@/components/modals/NewPermissionModal.vue'
import BTable, { BTableRow, BTableColumn } from '@/components/bootstrap/BTable.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    NewPermissionModal,
    AuthorityStatusBox,
  },
  setup() {
    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

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
      permissions,
      created,
    }
  },
})
</script>

