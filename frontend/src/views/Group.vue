<template>
  <div>
    <div class="container information">
      <h2>Group Information</h2>
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
import { ref, computed, onMounted, defineComponent } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { Group, Role } from '@/types/api'

import AttachRoleModal from '@/components/modals/AttachRoleModal.vue'
import DestructionModal from '@/components/modals/DestructionModal.vue'
import BButton from '@/components/bootstrap/BButton.vue'
import BTable, { BTableRow } from '@/components/bootstrap/BTable.vue'

export default defineComponent({
  components: {
    BButton,
    BTable,
    AttachRoleModal,
    DestructionModal,
  },
  setup() {
    const router = useRouter()
    const id = router.currentRoute.value.params['id'] as string

    const group = ref<Group>({
      id: '',
      name: '',
      description: '',
      attached_roles: [],
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
      } else {
        basicRow.value = []
      }
    })

    onMounted(() => {
      fetchGroup()
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
      group,
      basicRow,
      basicColumns,
      attachedRoles,
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