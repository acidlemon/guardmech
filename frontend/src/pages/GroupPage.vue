<template>
  <div>
    <div class="container information">
      <h2>Group Information</h2>
      <AuthorityStatusBox :status="authorityStatus" />
      <BTable
        v-if="group"
        :data="basicRow"
        :columns="basicColumns"
      />

      <h3>Attached Roles</h3>
      <AttachRoleModal
        owner-type="group"
        :owner-id="group.id"
        :attached-roles="attachedRoles"
        @completed="needRefresh"
      />
      <BTable
        v-if="attachedRoles.length"
        :data="attachedRoles"
        :columns="standardColumns"
      >
        <template #cell(action)="data">
          <!-- TODO: confirm modal -->
          <BButton variant="danger" @click="onDetachRole(data.row.id)">Detach</BButton>
        </template>
      </BTable>
      <p v-else>No roles.</p>

      <h3>Having Permissions</h3>
      <BTable
        v-if="permissionTableRows.length"
        :data="permissionTableRows"
        :columns="standardColumns"
      >
        <template #cell(action)="data">
          <template v-if="data && data.row.from">
            from {{ data.row.from }}
          </template>
        </template>
      </BTable>
      <p v-else>No permissions.</p>
    </div>

    <template v-if="group">
      <div class="danger-zone">
        <div class="container">
          <h3>Danger Zone</h3>
          <DestructionModal
            button-title="Delete This Group"
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
import { Group, Role } from '@/types/api'
import { useUserAuthority } from '@/hooks/useUserAuthority'

import AttachRoleModal from '@/components/modals/AttachRoleModal.vue'
import DestructionModal from '@/components/modals/DestructionModal.vue'
import BButton from '@/components/bootstrap/BButton.vue'
import BTable from '@/components/bootstrap/BTable.vue'
import AuthorityStatusBox from '@/components/AuthorityStatusBox.vue'
import { BTableRow } from '@/types/bootstrap'

export default defineComponent({
  components: {
    BButton,
    BTable,
    AttachRoleModal,
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

    const group = ref<Group>({
      id: '',
      name: '',
      description: '',
      attached_roles: [],
      having_permissions: [],
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
    const attachedRoles = computed<Role[]>(() => group.value ? group.value.attached_roles : [])

    const permissionTableRows = ref<BTableRow[]>([])

    const fetchGroup = (async () => {
      const res = await axios.get('/api/group/' + id)
      group.value = res.data.group as Group 

      if (group.value) {
        basicRow.value = [{
          label: 'Name',
          value: group.value.name,
        }, {
          label: 'Description',
          value: group.value.description,
        }]

        permissionTableRows.value = group.value.having_permissions.map(x => {
        const roles = group.value.attached_roles.filter(r => r.attached_permissions.find(p => p.id === x.id) ? true : false)

        return {
          ...x,
          from: roles.map(r => r.name).join(),
        }
      })
      } else {
        basicRow.value = []
      }
    })

    watch(authorityLoadCompleted, (val) => {
      if (val) {
        if (canRead.value) {
          fetchGroup()
        }
      }
    })

    const needRefresh = (() => {
      fetchGroup()
    })

    const onDetachRole = ((roleId: string) => {
      detachRole(roleId)
    })

    const detachRole = (async (roleId: string) => {
      const params = new URLSearchParams({
        role_id: roleId,
      })
      const res = await axios.post('/api/group/' + id + '/detach_role', params)
      console.log(res)

      group.value = res.data.group
      //fetchPrincipal()
    })

    const onDelete = (() => {
      deleteGroup()
    })
    const deleteGroup = (async () => {
      const res = await axios.delete('/api/group/' + id)
      console.log(res)
      // TODO check
      
      router.push('/groups')
    })

    return {
      authorityStatus,
      canWrite,
      canRead,
      group,
      basicRow,
      basicColumns,
      attachedRoles,
      permissionTableRows,
      standardColumns,
      needRefresh,
      onDelete,
      onDetachRole,
    }
  },
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