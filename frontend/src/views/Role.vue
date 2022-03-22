<template>
  <div>
    <div class="container information">
      <h2>Role Information</h2>
      <AuthorityStatusBox :status="authorityStatus" />
      <BTable
        v-if="role"
        :data="basicRow"
        :columns="basicColumns"
      />

      <h3>Attached Permissions</h3>
      <AttachPermissionModal
        :owner-id="role.id"
        :attached-permissions="attachedPermissions"
        @completed="needRefresh"
      />
      <BTable
        v-if="attachedPermissions.length"
        :data="attachedPermissions"
        :columns="standardColumns"
      >
        <template #cell(action)="data">
          <!-- TODO: confirm modal -->
          <BButton variant="danger" @click="onDetachPermission(data.row.id)">Detach</BButton>
        </template>
      </BTable>
      <p v-else>No attached permissions.</p>
    </div>

    <template v-if="role && canWrite">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Role"
            @confirmDelete="onDelete"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { ref, computed, watch, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Role, Permission } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import AttachPermissionModal from '@/components/modals/AttachPermissionModal.vue'
import DestructionModal from '@/components/modals/DestructionModal.vue'
import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow } from '@/components/bootstrap/BTable.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    AttachPermissionModal,
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

    const role = ref<Role>({
      id: '',
      name: '',
      description: '',
      attached_permissions: [],
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

    const standardColumns = [
      {
        key: 'name',
        label: 'Name',
      },
      {
        key: 'description',
        label: 'Description',
      },
      {
        key: 'action',
        label: '',
      },
    ]
    const attachedPermissions = computed<Permission[]>(() => role.value ? role.value.attached_permissions : [])

    const fetchRole = (async () => {
      const res = await axios.get('/api/role/' + id)
      role.value = res.data.role as Role

      if (role.value) {
        basicRow.value = [{
          label: 'Name',
          value: role.value.name,
        }, {
          label: 'Description',
          value: role.value.description,
        }]
      } else {
        basicRow.value = []
      }
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchRole()
        }
      }
    })

    const needRefresh = (() => {
      fetchRole()
    })

    const onDetachPermission = ((permissionId: string) => {
      detachPermission(permissionId)
    })

    const detachPermission = (async (permissionId: string) => {
      const params = new URLSearchParams({
        permission_id: permissionId,
      })
      const res = await axios.post('/api/role/' + id + '/detach_permission', params)
      console.log(res)
      fetchRole()
    })

    const onDelete = (() => {
      deleteRole()
    })
    const deleteRole = (async () => {
      const res = await axios.delete('/api/role/' + id)
      console.log(res)
      // TODO check
      
      router.push('/roles')
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      role,
      basicRow,
      basicColumns,
      standardColumns,
      attachedPermissions,
      needRefresh,
      onDelete,
      onDetachPermission,
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