<template>
  <div>
    <div class="container information">
      <h2>Permission Information</h2>
      <AuthorityStatusBox :status="authorityStatus" />
      <BTable
        v-if="permission"
        :data="basicRow"
        :columns="basicColumns"
      />
    </div>
    <template v-if="permission">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Permission"
            @confirmDelete="onDelete"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { ref, watch, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Permission } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import DestructionModal from '@/components/modals/DestructionModal.vue'
// import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow } from '@/types/bootstrap'

export default defineComponent({
  components: {
    // BButton,
    BTable,
    DestructionModal,
    AuthorityStatusBox,
  },
  setup() {
    const router = useRouter()
    const id = router.currentRoute.value.params['id'] as string

    const {
      authorityStatus,
      authorityLoadCompleted,
      canWrite,
      canRead,
    } = useUserAuthority()

    const permission = ref<Permission>({
      id: '',
      name: '',
      description: '',
    })

    const basicRow = ref<BTableRow[]>([])
    const basicColumns = [
      {
        key: 'label',
        label: 'Item',
      },
      {
        key: 'value',
        label: 'Data',
      },
    ]

    const fetchPermission = (async () => {
      const res = await axios.get('/api/permission/' + id)
      permission.value = res.data.permission as Permission

      if (permission.value) {
        basicRow.value = [{
          label: 'Name',
          value: permission.value.name,
        }, {
          label: 'Description',
          value: permission.value.description,
        }]
      } else {
        basicRow.value = []
      }
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchPermission()
        }
      }
    })

    const onDelete = (() => {
      deletePermission()
    })
    const deletePermission = (async () => {
      const res = await axios.delete('/api/permission/' + id)
      console.log(res)
      // TODO check

      router.push('/permissions')
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      permission,
      basicRow,
      basicColumns,
      onDelete,
    }
  }
})
</script>

<style scoped>
.information {
  padding-top: 20px;
  padding-bottom: 20px;
}

.danger-zone {
  background: #F8E0E0;
  padding-top: 20px;
  padding-bottom: 30px;
  border-top: dashed 1px #E06060;
}
</style>